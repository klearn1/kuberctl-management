/*
Copyright 2014 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package proxy

import (
	"fmt"
	"io"
	"net"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/golang/glog"
)

type serviceInfo struct {
	ip       net.IP
	port     int
	listener net.Listener
	mu       sync.Mutex // protects active
	active   bool
}

// Proxier is a simple proxy for TCP connections between a localhost:lport
// and services that provide the actual implementations.
type Proxier struct {
	loadBalancer LoadBalancer
	mu           sync.Mutex // protects serviceMap
	serviceMap   map[string]*serviceInfo
}

// NewProxier returns a new Proxier given a LoadBalancer.
func NewProxier(loadBalancer LoadBalancer) *Proxier {
	return &Proxier{
		loadBalancer: loadBalancer,
		serviceMap:   make(map[string]*serviceInfo),
	}
}

func copyBytes(in, out *net.TCPConn) {
	glog.Infof("Copying from %v <-> %v <-> %v <-> %v",
		in.RemoteAddr(), in.LocalAddr(), out.LocalAddr(), out.RemoteAddr())
	if _, err := io.Copy(in, out); err != nil {
		glog.Errorf("I/O error: %v", err)
	}
	in.CloseRead()
	out.CloseWrite()
}

// proxyConnection proxies data bidirectionally between in and out.
func proxyConnection(in, out *net.TCPConn) {
	glog.Infof("Creating proxy between %v <-> %v <-> %v <-> %v",
		in.RemoteAddr(), in.LocalAddr(), out.LocalAddr(), out.RemoteAddr())
	go copyBytes(in, out)
	go copyBytes(out, in)
}

// StopProxy stops the proxy for the named service.
func (proxier *Proxier) StopProxy(service string) error {
	// TODO: delete from map here?
	info, found := proxier.getServiceInfo(service)
	if !found {
		return fmt.Errorf("unknown service: %s", service)
	}
	return proxier.stopProxyInternal(service, info)
}

func (proxier *Proxier) stopProxyInternal(name string, info *serviceInfo) error {
	info.mu.Lock()
	defer info.mu.Unlock()
	if !info.active {
		return nil
	}
	glog.Infof("Removing service: %s", name)
	info.active = false
	return info.listener.Close()
}

func (proxier *Proxier) getServiceInfo(service string) (*serviceInfo, bool) {
	proxier.mu.Lock()
	defer proxier.mu.Unlock()
	info, ok := proxier.serviceMap[service]
	return info, ok
}

func (proxier *Proxier) setServiceInfo(service string, info *serviceInfo) {
	proxier.mu.Lock()
	defer proxier.mu.Unlock()
	proxier.serviceMap[service] = info
}

// AcceptHandler proxies incoming connections for the specified service
// to the load-balanced service endpoints.
func (proxier *Proxier) AcceptHandler(service string, listener net.Listener) {
	info, found := proxier.getServiceInfo(service)
	if !found {
		glog.Errorf("Failed to find service: %s", service)
		return
	}
	for {
		info.mu.Lock()
		if !info.active {
			info.mu.Unlock()
			glog.Infof("Cancelling Accept() loop for service: %s", service)
			break
		}
		info.mu.Unlock()
		inConn, err := listener.Accept()
		if err != nil {
			glog.Errorf("Accept failed: %v", err)
			continue
		}
		glog.Infof("Accepted connection from: %v to %v", inConn.RemoteAddr(), inConn.LocalAddr())
		endpoint, err := proxier.loadBalancer.NextEndpoint(service, inConn.RemoteAddr())
		if err != nil {
			glog.Errorf("Couldn't find an endpoint for %s %v", service, err)
			inConn.Close()
			continue
		}
		glog.Infof("Mapped service %s to endpoint %s", service, endpoint)
		outConn, err := net.DialTimeout("tcp", endpoint, time.Duration(5)*time.Second)
		if err != nil {
			glog.Errorf("Dial failed: %v", err)
			inConn.Close()
			continue
		}
		proxyConnection(inConn.(*net.TCPConn), outConn.(*net.TCPConn))
	}
}

// getListener decides which local port to listen on and returns a Listener and
// the assigned port.
func getListener(ip net.IP, port int) (net.Listener, int, error) {
	// If the portal IP is set, allocate a random port locally.
	if len(ip) != 0 {
		port = 0
	}
	//FIXME: if the portal IP is set, listen on localhost only?
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, 0, err
	}
	_, assignedPortStr, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		l.Close()
		return nil, 0, err
	}
	assignedPortNum, err := strconv.Atoi(assignedPortStr)
	if err != nil {
		l.Close()
		return nil, 0, err
	}
	return l, assignedPortNum, nil
}

// used to globally lock around unused ports. Only used in testing.
var unusedPortLock sync.Mutex

// addServiceOnUnusedPort starts listening for a new service, returning the
// port it's using.  For testing on a system with unknown ports used.
// FIXME: remove this?
func (proxier *Proxier) addServiceOnUnusedPort(service string) (string, error) {
	unusedPortLock.Lock()
	defer unusedPortLock.Unlock()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}
	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return "", err
	}
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return "", err
	}
	proxier.setServiceInfo(service, &serviceInfo{
		port:     portNum,
		active:   true,
		listener: l,
	})
	proxier.startAccepting(service, l)
	return port, nil
}

func (proxier *Proxier) startAccepting(service string, l net.Listener) {
	glog.Infof("Listening for %s on %s", service, l.Addr().String())
	go proxier.AcceptHandler(service, l)
}

func installRedirect(ip net.IP, port int, localPort int) error {
	const iptables = "iptables"

	//FIXME: check with -C first
	glog.Info("TIM: starting iptables rules")
	rule1 := []string{
		"-t", "nat",
		"-A", "PREROUTING",
		"-p", "tcp",
		"-d", ip.String(),
		"--dport", fmt.Sprintf("%d", port),
		"-j", "REDIRECT",
		"--to-ports", fmt.Sprintf("%d", localPort),
	}
	out, err := exec.Command(iptables, rule1...).CombinedOutput()
	if err != nil {
		glog.Errorf("error on iptables rule1: %s", out)
		return err
	}

	rule2 := []string{
		"-t", "nat",
		"-A", "OUTPUT",
		"-p", "tcp",
		"-d", ip.String(),
		"--dport", fmt.Sprintf("%d", port),
		"-j", "REDIRECT",
		"--to-ports", fmt.Sprintf("%d", localPort),
	}
	out, err = exec.Command(iptables, rule2...).CombinedOutput()
	if err != nil {
		glog.Errorf("error on iptables rule2: %s", out)
		//FIXME: undo rule1
		return err
	}

	glog.Info("TIM: iptables rules are installed")
	return nil
}

// OnUpdate manages the active set of service proxies.
// Active service proxies are reinitialized if found in the update set or
// shutdown if missing from the update set.
func (proxier *Proxier) OnUpdate(services []api.Service) {
	glog.Infof("Received update notice: %+v", services)
	activeServices := util.StringSet{}
	for _, service := range services {
		activeServices.Insert(service.ID)
		info, exists := proxier.getServiceInfo(service.ID)
		if exists && info.active && info.port == service.Port {
			continue
		}
		//FIXME: also handle the IP changing
		if exists && info.port != service.Port {
			//FIXME: remove the iptables rules
			proxier.StopProxy(service.ID)
		}
		glog.Infof("Adding a new service %s at %s:%d", service.ID, service.PortalIP, service.Port)
		portalIP := net.ParseIP(service.PortalIP)
		listener, port, err := getListener(portalIP, service.Port)
		if err != nil {
			glog.Infof("Failed to start listening for %s: %+v", service.ID, err)
			continue
		}
		glog.Infof("Proxying for service %s on local port %d", service.ID, port)
		err = installRedirect(portalIP, service.Port, port)
		if err != nil {
			glog.Infof("Failed to create IP redirect for %s on %d", service.ID, service.Port)
			listener.Close()
			continue
		}
		proxier.setServiceInfo(service.ID, &serviceInfo{
			ip:       portalIP,
			port:     service.Port,
			active:   true,
			listener: listener,
		})
		proxier.startAccepting(service.ID, listener)
	}
	proxier.mu.Lock()
	defer proxier.mu.Unlock()
	for name, info := range proxier.serviceMap {
		if !activeServices.Has(name) {
			//FIXME: remove the iptables rules
			proxier.stopProxyInternal(name, info)
		}
	}
}
