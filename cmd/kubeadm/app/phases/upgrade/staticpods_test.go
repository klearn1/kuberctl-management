/*
Copyright 2017 The Kubernetes Authors.

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

package upgrade

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/pkg/errors"

	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
	certsphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
	controlplanephase "k8s.io/kubernetes/cmd/kubeadm/app/phases/controlplane"
	etcdphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/etcd"
	"k8s.io/kubernetes/cmd/kubeadm/app/util"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/apiclient"
	certstestutil "k8s.io/kubernetes/cmd/kubeadm/app/util/certs"
	configutil "k8s.io/kubernetes/cmd/kubeadm/app/util/config"
	etcdutil "k8s.io/kubernetes/cmd/kubeadm/app/util/etcd"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/pkiutil"
	testutil "k8s.io/kubernetes/cmd/kubeadm/test"
)

const (
	waitForHashes         = "wait-for-hashes"
	waitForHashChange     = "wait-for-hash-change"
	waitForPodsWithLabel  = "wait-for-pods-with-label"
	clusterStatusEndpoint = "https://192.168.2.2:2379"

	testConfiguration = `
apiVersion: kubeadm.k8s.io/v1beta1
kind: InitConfiguration
nodeRegistration:
  name: foo
  criSocket: ""
localAPIEndpoint:
  advertiseAddress: 192.168.2.2
  bindPort: 6443
bootstrapTokens:
- token: ce3aa5.5ec8455bb76b379f
  ttl: 24h
---
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration

apiServer:
  certSANs: null
  extraArgs: null
certificatesDir: %s
etcd:
  local:
    dataDir: %s
    image: ""
imageRepository: k8s.gcr.io
kubernetesVersion: %s
networking:
  dnsDomain: cluster.local
  podSubnet: ""
  serviceSubnet: 10.96.0.0/12
useHyperKubeImage: false
`
)

// fakeWaiter is a fake apiclient.Waiter that returns errors it was initialized with
type fakeWaiter struct {
	errsToReturn map[string]error
}

func NewFakeStaticPodWaiter(errsToReturn map[string]error) apiclient.Waiter {
	return &fakeWaiter{
		errsToReturn: errsToReturn,
	}
}

// WaitForAPI just returns a dummy nil, to indicate that the program should just proceed
func (w *fakeWaiter) WaitForAPI() error {
	return nil
}

// WaitForPodsWithLabel just returns an error if set from errsToReturn
func (w *fakeWaiter) WaitForPodsWithLabel(kvLabel string) error {
	return w.errsToReturn[waitForPodsWithLabel]
}

// WaitForPodToDisappear just returns a dummy nil, to indicate that the program should just proceed
func (w *fakeWaiter) WaitForPodToDisappear(podName string) error {
	return nil
}

// SetTimeout is a no-op; we don't use it in this implementation
func (w *fakeWaiter) SetTimeout(_ time.Duration) {}

// WaitForStaticPodControlPlaneHashes returns an error if set from errsToReturn
func (w *fakeWaiter) WaitForStaticPodControlPlaneHashes(_ string) (map[string]string, error) {
	return map[string]string{}, w.errsToReturn[waitForHashes]
}

// WaitForStaticPodSingleHash returns an error if set from errsToReturn
func (w *fakeWaiter) WaitForStaticPodSingleHash(_ string, _ string) (string, error) {
	return "", w.errsToReturn[waitForHashes]
}

// WaitForStaticPodHashChange returns an error if set from errsToReturn
func (w *fakeWaiter) WaitForStaticPodHashChange(_, _, _ string) error {
	return w.errsToReturn[waitForHashChange]
}

// WaitForHealthyKubelet returns a dummy nil just to implement the interface
func (w *fakeWaiter) WaitForHealthyKubelet(_ time.Duration, _ string) error {
	return nil
}

// WaitForKubeletAndFunc is a wrapper for WaitForHealthyKubelet that also blocks for a function
func (w *fakeWaiter) WaitForKubeletAndFunc(f func() error) error {
	return nil
}

type fakeStaticPodPathManager struct {
	kubernetesDir     string
	realManifestDir   string
	tempManifestDir   string
	backupManifestDir string
	backupEtcdDir     string
	MoveFileFunc      func(string, string) error
}

func NewFakeStaticPodPathManager(moveFileFunc func(string, string) error) (StaticPodPathManager, error) {
	kubernetesDir, err := ioutil.TempDir("", "kubeadm-pathmanager-")
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't create a temporary directory for the upgrade")
	}

	realManifestDir := filepath.Join(kubernetesDir, constants.ManifestsSubDirName)
	if err := os.Mkdir(realManifestDir, 0700); err != nil {
		return nil, errors.Wrapf(err, "couldn't create a realManifestDir for the upgrade")
	}

	upgradedManifestDir := filepath.Join(kubernetesDir, "upgraded-manifests")
	if err := os.Mkdir(upgradedManifestDir, 0700); err != nil {
		return nil, errors.Wrapf(err, "couldn't create a upgradedManifestDir for the upgrade")
	}

	backupManifestDir := filepath.Join(kubernetesDir, "backup-manifests")
	if err := os.Mkdir(backupManifestDir, 0700); err != nil {
		return nil, errors.Wrap(err, "couldn't create a backupManifestDir for the upgrade")
	}

	backupEtcdDir := filepath.Join(kubernetesDir, "kubeadm-backup-etcd")
	if err := os.Mkdir(backupEtcdDir, 0700); err != nil {
		return nil, err
	}

	return &fakeStaticPodPathManager{
		kubernetesDir:     kubernetesDir,
		realManifestDir:   realManifestDir,
		tempManifestDir:   upgradedManifestDir,
		backupManifestDir: backupManifestDir,
		backupEtcdDir:     backupEtcdDir,
		MoveFileFunc:      moveFileFunc,
	}, nil
}

func (spm *fakeStaticPodPathManager) MoveFile(oldPath, newPath string) error {
	return spm.MoveFileFunc(oldPath, newPath)
}

func (spm *fakeStaticPodPathManager) KubernetesDir() string {
	return spm.kubernetesDir
}

func (spm *fakeStaticPodPathManager) RealManifestPath(component string) string {
	return constants.GetStaticPodFilepath(component, spm.realManifestDir)
}
func (spm *fakeStaticPodPathManager) RealManifestDir() string {
	return spm.realManifestDir
}

func (spm *fakeStaticPodPathManager) TempManifestPath(component string) string {
	return constants.GetStaticPodFilepath(component, spm.tempManifestDir)
}
func (spm *fakeStaticPodPathManager) TempManifestDir() string {
	return spm.tempManifestDir
}

func (spm *fakeStaticPodPathManager) BackupManifestPath(component string) string {
	return constants.GetStaticPodFilepath(component, spm.backupManifestDir)
}
func (spm *fakeStaticPodPathManager) BackupManifestDir() string {
	return spm.backupManifestDir
}

func (spm *fakeStaticPodPathManager) BackupEtcdDir() string {
	return spm.backupEtcdDir
}

func (spm *fakeStaticPodPathManager) CleanupDirs() error {
	if err := os.RemoveAll(spm.TempManifestDir()); err != nil {
		return err
	}
	if err := os.RemoveAll(spm.BackupManifestDir()); err != nil {
		return err
	}
	return os.RemoveAll(spm.BackupEtcdDir())
}

func (spm *fakeStaticPodPathManager) CleanupKubernetesDir() error {
	return os.RemoveAll(spm.KubernetesDir())
}

type fakeTLSEtcdClient struct {
	TLS              bool
	version          string
	versionErr       bool
	clusterStatusErr bool
	availableErr     bool
}

func (c fakeTLSEtcdClient) ClusterAvailable() (bool, error) { return true, nil }

func (c fakeTLSEtcdClient) WaitForClusterAvailable(retries int, retryInterval time.Duration) (bool, error) {
	if c.availableErr {
		return false, errors.New("WaitForClusterAvailable failed")
	}
	return true, nil
}

func (c fakeTLSEtcdClient) GetClusterStatus() (map[string]*clientv3.StatusResponse, error) {
	if c.clusterStatusErr {
		return nil, errors.New("GetClusterStatus failed")
	}
	return map[string]*clientv3.StatusResponse{
		clusterStatusEndpoint: {
			Version: "3.1.12",
		}}, nil
}

func (c fakeTLSEtcdClient) GetClusterVersions() (map[string]string, error) {
	if c.versionErr {
		return nil, errors.New("Unsupported or unknown Kubernetes version")
	}
	if c.version != "" {
		return map[string]string{
			clusterStatusEndpoint: c.version,
		}, nil
	}
	return map[string]string{
		clusterStatusEndpoint: "3.1.12",
	}, nil
}

func (c fakeTLSEtcdClient) GetVersion() (string, error) {
	return "3.1.12", nil
}

func (c fakeTLSEtcdClient) Sync() error { return nil }

func (c fakeTLSEtcdClient) AddMember(name string, peerAddrs string) ([]etcdutil.Member, error) {
	return []etcdutil.Member{}, nil
}

func (c fakeTLSEtcdClient) GetMemberID(peerURL string) (uint64, error) {
	return 0, nil
}

func (c fakeTLSEtcdClient) RemoveMember(id uint64) ([]etcdutil.Member, error) {
	return []etcdutil.Member{}, nil
}

type fakePodManifestEtcdClient struct {
	ManifestDir, CertificatesDir string
	AvailableErr                 bool
}

func (c fakePodManifestEtcdClient) ClusterAvailable() (bool, error) { return true, nil }

func (c fakePodManifestEtcdClient) WaitForClusterAvailable(retries int, retryInterval time.Duration) (bool, error) {
	if c.AvailableErr {
		return false, errors.New("WaitForClusterAvailable failed")
	}
	return true, nil
}

func (c fakePodManifestEtcdClient) GetClusterStatus() (map[string]*clientv3.StatusResponse, error) {
	// Make sure the certificates generated from the upgrade are readable from disk
	tlsInfo := transport.TLSInfo{
		CertFile:      filepath.Join(c.CertificatesDir, constants.EtcdCACertName),
		KeyFile:       filepath.Join(c.CertificatesDir, constants.EtcdHealthcheckClientCertName),
		TrustedCAFile: filepath.Join(c.CertificatesDir, constants.EtcdHealthcheckClientKeyName),
	}
	_, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, err
	}

	return map[string]*clientv3.StatusResponse{
		clusterStatusEndpoint: {Version: "3.1.12"},
	}, nil
}

func (c fakePodManifestEtcdClient) GetClusterVersions() (map[string]string, error) {
	return map[string]string{
		clusterStatusEndpoint: "3.1.12",
	}, nil
}

func (c fakePodManifestEtcdClient) GetVersion() (string, error) {
	return "3.1.12", nil
}

func (c fakePodManifestEtcdClient) Sync() error { return nil }

func (c fakePodManifestEtcdClient) AddMember(name string, peerAddrs string) ([]etcdutil.Member, error) {
	return []etcdutil.Member{}, nil
}

func (c fakePodManifestEtcdClient) GetMemberID(peerURL string) (uint64, error) {
	return 0, nil
}

func (c fakePodManifestEtcdClient) RemoveMember(id uint64) ([]etcdutil.Member, error) {
	return []etcdutil.Member{}, nil
}

func TestStaticPodControlPlane(t *testing.T) {
	tests := []struct {
		description          string
		waitErrsToReturn     map[string]error
		moveFileFunc         func(string, string) error
		expectedErr          bool
		manifestShouldChange bool
	}{
		{
			description: "error-free case should succeed",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          false,
			manifestShouldChange: true,
		},
		{
			description: "any wait error should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        errors.New("boo! failed"),
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
		{
			description: "any wait error should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    errors.New("boo! failed"),
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
		{
			description: "any wait error should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: errors.New("boo! failed"),
			},
			moveFileFunc: func(oldPath, newPath string) error {
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
		{
			description: "any path-moving error should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				// fail for kube-apiserver move
				if strings.Contains(newPath, "kube-apiserver") {
					return errors.New("moving the kube-apiserver file failed")
				}
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
		{
			description: "any path-moving error should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				// fail for kube-controller-manager move
				if strings.Contains(newPath, "kube-controller-manager") {
					return errors.New("moving the kube-apiserver file failed")
				}
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
		{
			description: "any path-moving error should result in a rollback and an abort; even though this is the last component (kube-apiserver and kube-controller-manager healthy)",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			moveFileFunc: func(oldPath, newPath string) error {
				// fail for kube-scheduler move
				if strings.Contains(newPath, "kube-scheduler") {
					return errors.New("moving the kube-apiserver file failed")
				}
				return os.Rename(oldPath, newPath)
			},
			expectedErr:          true,
			manifestShouldChange: false,
		},
	}

	for _, rt := range tests {
		t.Run(rt.description, func(t *testing.T) {
			waiter := NewFakeStaticPodWaiter(rt.waitErrsToReturn)
			pathMgr, err := NewFakeStaticPodPathManager(rt.moveFileFunc)
			if err != nil {
				t.Fatalf("couldn't run NewFakeStaticPodPathManager: %v", err)
			}
			defer os.RemoveAll(pathMgr.(*fakeStaticPodPathManager).KubernetesDir())
			constants.KubernetesDir = pathMgr.(*fakeStaticPodPathManager).KubernetesDir()

			tempCertsDir, err := ioutil.TempDir("", "kubeadm-certs")
			if err != nil {
				t.Fatalf("couldn't create temporary certificates directory: %v", err)
			}
			defer os.RemoveAll(tempCertsDir)
			tmpEtcdDataDir, err := ioutil.TempDir("", "kubeadm-etcd-data")
			if err != nil {
				t.Fatalf("couldn't create temporary etcd data directory: %v", err)
			}
			defer os.RemoveAll(tmpEtcdDataDir)

			oldcfg, err := getConfig(constants.MinimumControlPlaneVersion.String(), tempCertsDir, tmpEtcdDataDir)
			if err != nil {
				t.Fatalf("couldn't create config: %v", err)
			}

			tree, err := certsphase.GetCertsWithoutEtcd().AsMap().CertTree()
			if err != nil {
				t.Fatalf("couldn't get cert tree: %v", err)
			}

			if err := tree.CreateTree(oldcfg); err != nil {
				t.Fatalf("couldn't get create cert tree: %v", err)
			}

			t.Logf("Wrote certs to %s\n", oldcfg.CertificatesDir)

			// Initialize the directory with v1.7 manifests; should then be upgraded to v1.8 using the method
			err = controlplanephase.CreateInitStaticPodManifestFiles(pathMgr.RealManifestDir(), oldcfg)
			if err != nil {
				t.Fatalf("couldn't run CreateInitStaticPodManifestFiles: %v", err)
			}
			err = etcdphase.CreateLocalEtcdStaticPodManifestFile(pathMgr.RealManifestDir(), oldcfg.NodeRegistration.Name, &oldcfg.ClusterConfiguration, &oldcfg.LocalAPIEndpoint)
			if err != nil {
				t.Fatalf("couldn't run CreateLocalEtcdStaticPodManifestFile: %v", err)
			}
			// Get a hash of the v1.7 API server manifest to compare later (was the file re-written)
			oldHash, err := getAPIServerHash(pathMgr.RealManifestDir())
			if err != nil {
				t.Fatalf("couldn't read temp file: %v", err)
			}

			newcfg, err := getConfig(constants.CurrentKubernetesVersion.String(), tempCertsDir, tmpEtcdDataDir)
			if err != nil {
				t.Fatalf("couldn't create config: %v", err)
			}

			// create the kubeadm etcd certs
			caCert, caKey, err := certsphase.KubeadmCertEtcdCA.CreateAsCA(newcfg)
			if err != nil {
				t.Fatalf("couldn't create new CA certificate: %v", err)
			}
			for _, cert := range []*certsphase.KubeadmCert{
				&certsphase.KubeadmCertEtcdServer,
				&certsphase.KubeadmCertEtcdPeer,
				&certsphase.KubeadmCertEtcdHealthcheck,
				&certsphase.KubeadmCertEtcdAPIClient,
			} {
				if err := cert.CreateFromCA(newcfg, caCert, caKey); err != nil {
					t.Fatalf("couldn't create certificate %s: %v", cert.Name, err)
				}
			}

			actualErr := StaticPodControlPlane(
				nil,
				waiter,
				pathMgr,
				newcfg,
				true,
				fakeTLSEtcdClient{
					TLS: false,
				},
				fakePodManifestEtcdClient{
					ManifestDir:     pathMgr.RealManifestDir(),
					CertificatesDir: newcfg.CertificatesDir,
				},
			)
			if (actualErr != nil) != rt.expectedErr {
				t.Errorf(
					"failed UpgradeStaticPodControlPlane\n%s\n\texpected error: %t\n\tgot: %t\n\tactual error: %v",
					rt.description,
					rt.expectedErr,
					(actualErr != nil),
					actualErr,
				)
			}

			newHash, err := getAPIServerHash(pathMgr.RealManifestDir())
			if err != nil {
				t.Fatalf("couldn't read temp file: %v", err)
			}

			if (oldHash != newHash) != rt.manifestShouldChange {
				t.Errorf(
					"failed StaticPodControlPlane\n%s\n\texpected manifest change: %t\n\tgot: %t\n\tnewHash: %v",
					rt.description,
					rt.manifestShouldChange,
					(oldHash != newHash),
					newHash,
				)
			}
		})
	}
}

func getAPIServerHash(dir string) (string, error) {
	manifestPath := constants.GetStaticPodFilepath(constants.KubeAPIServer, dir)

	fileBytes, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(fileBytes)), nil
}

func getConfig(version, certsDir, etcdDataDir string) (*kubeadmapi.InitConfiguration, error) {
	configBytes := []byte(fmt.Sprintf(testConfiguration, certsDir, etcdDataDir, version))

	// Unmarshal the config
	return configutil.BytesToInitConfiguration(configBytes)
}

func getTempDir(t *testing.T, name string) (string, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), name)
	if err != nil {
		t.Fatalf("couldn't make temporary directory: %v", err)
	}

	return dir, func() {
		os.RemoveAll(dir)
	}
}

func getPathMgrAndCfg(t *testing.T) (StaticPodPathManager, *kubeadmapi.InitConfiguration) {
	pathMgr, err := NewFakeStaticPodPathManager(os.Rename)
	if err != nil {
		t.Fatalf("couldn't run NewFakeStaticPodPathManager: %v", err)
	}

	// constants.KubernetesDir = pathMgr.(*fakeStaticPodPathManager).KubernetesDir()

	tempCertsDir, err := ioutil.TempDir("", "kubeadm-certs")
	if err != nil {
		t.Fatalf("couldn't create temporary certificates directory: %v", err)
	}

	tmpEtcdDataDir, err := ioutil.TempDir("", "kubeadm-etcd-data")
	if err != nil {
		t.Fatalf("couldn't create temporary etcd data directory: %v", err)
	}

	oldcfg, err := getConfig(constants.MinimumControlPlaneVersion.String(), tempCertsDir, tmpEtcdDataDir)
	if err != nil {
		t.Fatalf("couldn't create config: %v", err)
	}

	tree, err := certsphase.GetCertsWithoutEtcd().AsMap().CertTree()
	if err != nil {
		t.Fatalf("couldn't get cert tree: %v", err)
	}

	if err := tree.CreateTree(oldcfg); err != nil {
		t.Fatalf("couldn't get create cert tree: %v", err)
	}

	t.Logf("Wrote certs to %s\n", oldcfg.CertificatesDir)

	// Initialize the directory with v1.13 manifests; should then be upgraded to v1.14 using the method
	err = controlplanephase.CreateInitStaticPodManifestFiles(pathMgr.RealManifestDir(), oldcfg)
	if err != nil {
		t.Fatalf("couldn't run CreateInitStaticPodManifestFiles: %v", err)
	}
	err = etcdphase.CreateLocalEtcdStaticPodManifestFile(pathMgr.RealManifestDir(), oldcfg.NodeRegistration.Name, &oldcfg.ClusterConfiguration, &oldcfg.LocalAPIEndpoint)
	if err != nil {
		t.Fatalf("couldn't run CreateLocalEtcdStaticPodManifestFile: %v", err)
	}

	newcfg, err := getConfig(constants.CurrentKubernetesVersion.String(), tempCertsDir, tmpEtcdDataDir)
	if err != nil {
		t.Fatalf("couldn't create config: %v", err)
	}

	// create the kubeadm etcd certs
	caCert, caKey, err := certsphase.KubeadmCertEtcdCA.CreateAsCA(newcfg)
	if err != nil {
		t.Fatalf("couldn't create new CA certificate: %v", err)
	}
	for _, cert := range []*certsphase.KubeadmCert{
		&certsphase.KubeadmCertEtcdServer,
		&certsphase.KubeadmCertEtcdPeer,
		&certsphase.KubeadmCertEtcdHealthcheck,
		&certsphase.KubeadmCertEtcdAPIClient,
	} {
		if err := cert.CreateFromCA(newcfg, caCert, caKey); err != nil {
			t.Fatalf("couldn't create certificate %s: %v", cert.Name, err)
		}
	}
	return pathMgr, newcfg
}

func TestCleanupDirs(t *testing.T) {
	tests := []struct {
		name                   string
		keepManifest, keepEtcd bool
	}{
		{
			name:         "save manifest backup",
			keepManifest: true,
		},
		{
			name:         "save both etcd and manifest",
			keepManifest: true,
			keepEtcd:     true,
		},
		{
			name: "save nothing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			realManifestDir, cleanup := getTempDir(t, "realManifestDir")
			defer cleanup()

			tempManifestDir, cleanup := getTempDir(t, "tempManifestDir")
			defer cleanup()

			backupManifestDir, cleanup := getTempDir(t, "backupManifestDir")
			defer cleanup()

			backupEtcdDir, cleanup := getTempDir(t, "backupEtcdDir")
			defer cleanup()

			mgr := NewKubeStaticPodPathManager(realManifestDir, tempManifestDir, backupManifestDir, backupEtcdDir, test.keepManifest, test.keepEtcd)
			err := mgr.CleanupDirs()
			if err != nil {
				t.Errorf("unexpected error cleaning up: %v", err)
			}

			if _, err := os.Stat(tempManifestDir); !os.IsNotExist(err) {
				t.Errorf("%q should not have existed", tempManifestDir)
			}
			_, err = os.Stat(backupManifestDir)
			if test.keepManifest {
				if err != nil {
					t.Errorf("unexpected error getting backup manifest dir")
				}
			} else {
				if !os.IsNotExist(err) {
					t.Error("expected backup manifest to not exist")
				}
			}

			_, err = os.Stat(backupEtcdDir)
			if test.keepEtcd {
				if err != nil {
					t.Errorf("unexpected error getting backup etcd dir")
				}
			} else {
				if !os.IsNotExist(err) {
					t.Error("expected backup etcd dir to not exist")
				}
			}
		})
	}
}

func TestRenewCerts(t *testing.T) {
	caCert, caKey := certstestutil.SetupCertificateAuthorithy(t)
	t.Run("all certs exist, should be rotated", func(t *testing.T) {
	})
	tests := []struct {
		name               string
		component          string
		skipCreateCA       bool
		shouldErrorOnRenew bool
		certsShouldExist   []*certsphase.KubeadmCert
	}{
		{
			name:      "all certs exist, should be rotated",
			component: constants.Etcd,
			certsShouldExist: []*certsphase.KubeadmCert{
				&certsphase.KubeadmCertEtcdServer,
				&certsphase.KubeadmCertEtcdPeer,
				&certsphase.KubeadmCertEtcdHealthcheck,
			},
		},
		{
			name:      "just renew API cert",
			component: constants.KubeAPIServer,
			certsShouldExist: []*certsphase.KubeadmCert{
				&certsphase.KubeadmCertEtcdAPIClient,
			},
		},
		{
			name:         "ignores other compnonents",
			skipCreateCA: true,
			component:    constants.KubeScheduler,
		},
		{
			name:               "missing a cert to renew",
			component:          constants.Etcd,
			shouldErrorOnRenew: true,
			certsShouldExist: []*certsphase.KubeadmCert{
				&certsphase.KubeadmCertEtcdServer,
				&certsphase.KubeadmCertEtcdPeer,
			},
		},
		{
			name:               "no CA, cannot continue",
			component:          constants.Etcd,
			skipCreateCA:       true,
			shouldErrorOnRenew: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup up basic requities
			tmpDir := testutil.SetupTempDir(t)
			defer os.RemoveAll(tmpDir)

			cfg := testutil.GetDefaultInternalConfig(t)
			cfg.CertificatesDir = tmpDir

			if !test.skipCreateCA {
				if err := pkiutil.WriteCertAndKey(tmpDir, constants.EtcdCACertAndKeyBaseName, caCert, caKey); err != nil {
					t.Fatalf("couldn't write out CA: %v", err)
				}
			}

			// Create expected certs
			for _, kubeCert := range test.certsShouldExist {
				if err := kubeCert.CreateFromCA(cfg, caCert, caKey); err != nil {
					t.Fatalf("couldn't renew certificate %q: %v", kubeCert.Name, err)
				}
			}

			// Load expected certs to check if serial numbers changes
			certMaps := make(map[*certsphase.KubeadmCert]big.Int)
			for _, kubeCert := range test.certsShouldExist {
				cert, err := pkiutil.TryLoadCertFromDisk(tmpDir, kubeCert.BaseName)
				if err != nil {
					t.Fatalf("couldn't load certificate %q: %v", kubeCert.Name, err)
				}
				certMaps[kubeCert] = *cert.SerialNumber
			}

			// Renew everything
			err := renewCerts(cfg, test.component)
			if test.shouldErrorOnRenew {
				if err == nil {
					t.Fatal("expected renewal error, got nothing")
				}
				// expected error, got error
				return
			}
			if err != nil {
				t.Fatalf("couldn't renew certificates: %v", err)
			}

			// See if the certificate serial numbers change
			for kubeCert, cert := range certMaps {
				newCert, err := pkiutil.TryLoadCertFromDisk(tmpDir, kubeCert.BaseName)
				if err != nil {
					t.Errorf("couldn't load new certificate %q: %v", kubeCert.Name, err)
					continue
				}
				if cert.Cmp(newCert.SerialNumber) == 0 {
					t.Errorf("certifitate %v was not reissued", kubeCert.Name)
				}
			}
		})

	}
}

func TestCompareEtcdVersion(t *testing.T) {
	var tests = []struct {
		description   string
		cfg           *kubeadmapi.InitConfiguration
		oldEtcdClient etcdutil.ClusterInterrogator
		expectedFatal bool
		expectedError bool
	}{
		{
			description: "cfg.KubernetesVersion error",
			cfg: &kubeadmapi.InitConfiguration{
				ClusterConfiguration: kubeadmapi.ClusterConfiguration{
					KubernetesVersion: "1.99.0",
				},
			},
			expectedFatal: true,
			expectedError: true,
		},
		{
			description: "oldEtcdClient.GetClusterVersions error",
			cfg: &kubeadmapi.InitConfiguration{
				ClusterConfiguration: kubeadmapi.ClusterConfiguration{
					KubernetesVersion: "1.14.0",
				},
			},
			oldEtcdClient: fakeTLSEtcdClient{
				versionErr: true,
			},
			expectedFatal: true,
			expectedError: true,
		},
		{
			description: "cfg.LocalAPIEndpoint error",
			cfg: &kubeadmapi.InitConfiguration{
				LocalAPIEndpoint: kubeadmapi.APIEndpoint{AdvertiseAddress: "192.168.2.1"},
				ClusterConfiguration: kubeadmapi.ClusterConfiguration{
					KubernetesVersion: "1.14.0",
				},
			},
			oldEtcdClient: fakeTLSEtcdClient{},
			expectedFatal: true,
			expectedError: true,
		},
		{
			description: "parse etcdVersion error",
			cfg: &kubeadmapi.InitConfiguration{
				LocalAPIEndpoint: kubeadmapi.APIEndpoint{AdvertiseAddress: "192.168.2.2"},
				ClusterConfiguration: kubeadmapi.ClusterConfiguration{
					KubernetesVersion: "1.14.0",
				},
			},
			oldEtcdClient: fakeTLSEtcdClient{
				version: "3:1:22",
			},
			expectedFatal: true,
			expectedError: true,
		},
		{
			description: "desiredEtcdVersion lessThan currentEtcdVersion error",
			cfg: &kubeadmapi.InitConfiguration{
				LocalAPIEndpoint: kubeadmapi.APIEndpoint{AdvertiseAddress: "192.168.2.2"},
				ClusterConfiguration: kubeadmapi.ClusterConfiguration{
					KubernetesVersion: "1.13.1",
				},
			},
			oldEtcdClient: fakeTLSEtcdClient{
				version: "3.3.10",
			},
			expectedFatal: false,
			expectedError: true,
		},
		{
			description: "desiredEtcdVersion is the same as currentEtcdVersion",
			cfg: &kubeadmapi.InitConfiguration{
				LocalAPIEndpoint: kubeadmapi.APIEndpoint{AdvertiseAddress: "192.168.2.2"},
				ClusterConfiguration: kubeadmapi.ClusterConfiguration{
					KubernetesVersion: "1.14.0",
				},
			},
			oldEtcdClient: fakeTLSEtcdClient{
				version: "3.3.10",
			},
			expectedFatal: false,
			expectedError: false,
		},
		{
			description: "etcdVersion compare success",
			cfg: &kubeadmapi.InitConfiguration{
				LocalAPIEndpoint: kubeadmapi.APIEndpoint{AdvertiseAddress: "192.168.2.2"},
				ClusterConfiguration: kubeadmapi.ClusterConfiguration{
					KubernetesVersion: "1.14.0",
				},
			},
			oldEtcdClient: fakeTLSEtcdClient{
				version: "3.2.24",
			},
			expectedFatal: true,
			expectedError: false,
		},
	}

	for _, rt := range tests {
		t.Run(rt.description, func(t *testing.T) {
			actualFatal, actualError := compareEtcdVersion(rt.cfg, rt.oldEtcdClient)
			if actualFatal != rt.expectedFatal {
				t.Errorf("%s unexpected failure: %v", rt.description, actualFatal)
				return
			}
			if (actualError != nil) && !rt.expectedError {
				t.Errorf("%s unexpected failure: %v", rt.description, actualError)
				return
			} else if (actualError == nil) && rt.expectedError {
				t.Errorf("%s passed when expected to fail", rt.description)
				return
			}
		})
	}
}

func TestRollbackEtcdOnFailedUpgrade(t *testing.T) {
	tests := []struct {
		description      string
		oldEtcdClient    etcdutil.ClusterInterrogator
		recoverManifests map[string]string
		expectedErr      bool
	}{
		{
			description:   "roll back etcd success with no recoverManifests",
			oldEtcdClient: fakeTLSEtcdClient{},
			expectedErr:   false,
		},
		{
			description: "roll back etcd fail with no recoverManifests",
			oldEtcdClient: fakeTLSEtcdClient{
				availableErr: true,
			},
			expectedErr: true,
		},
		{
			description:      "roll back etcd success with recoverManifests",
			oldEtcdClient:    fakeTLSEtcdClient{},
			recoverManifests: map[string]string{},
			expectedErr:      false,
		},
		{
			description: "roll back etcd fail with recoverManifests",
			oldEtcdClient: fakeTLSEtcdClient{
				availableErr: true,
			},
			recoverManifests: map[string]string{},
			expectedErr:      true,
		},
	}

	for _, rt := range tests {
		pathMgr, newcfg := getPathMgrAndCfg(t)
		defer pathMgr.(*fakeStaticPodPathManager).CleanupKubernetesDir()
		defer os.RemoveAll(newcfg.CertificatesDir)
		defer os.RemoveAll(newcfg.Etcd.Local.DataDir)

		backupEtcdDir := pathMgr.BackupEtcdDir()
		runningEtcdDir := newcfg.Etcd.Local.DataDir
		if err := util.CopyDir(runningEtcdDir, backupEtcdDir); err != nil {
			t.Fatalf("failed to back up etcd data: %v", err)
		}

		var errForRecoverManifests error
		if rt.recoverManifests != nil {
			backupManifestPath := pathMgr.BackupManifestPath(constants.Etcd)
			rt.recoverManifests[constants.Etcd] = backupManifestPath
			errForRecoverManifests = errors.New("WaitForClusterAvailable failed")
		}

		actualErr := rollbackEtcdOnFailedUpgrade(
			pathMgr,
			newcfg,
			rt.recoverManifests,
			rt.oldEtcdClient,
			backupEtcdDir,
			errForRecoverManifests,
		)
		if (actualErr != nil) != rt.expectedErr {
			t.Errorf(
				"failed rollbackEtcdOnFailedUpgrade\n%s\n\texpected error: %t\n\tgot: %t\n\tactual error: %v",
				rt.description,
				rt.expectedErr,
				(actualErr != nil),
				actualErr,
			)
		}
	}
}

func TestPerformEtcdStaticPodUpgrade(t *testing.T) {
	tests := []struct {
		description      string
		waitErrsToReturn map[string]error
		oldEtcdClient    etcdutil.ClusterInterrogator
		newEtcdClient    etcdutil.ClusterInterrogator
		expectedFatal    bool
		expectedErr      bool
	}{
		{
			description: "error-free case should succeed",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			oldEtcdClient: fakeTLSEtcdClient{},
			expectedFatal: false,
			expectedErr:   false,
		},
		{
			description: "etcd cluster is not healthy should abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			oldEtcdClient: fakeTLSEtcdClient{
				clusterStatusErr: true,
			},
			expectedFatal: true,
			expectedErr:   true,
		},
		{
			description: "compare the currentEtcdVersion and desiredEtcdVersion failed should abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			oldEtcdClient: fakeTLSEtcdClient{
				version: "3:3:10",
			},
			expectedFatal: true,
			expectedErr:   true,
		},
		{
			description: "any wait error should abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        errors.New("boo! failed"),
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			oldEtcdClient: fakeTLSEtcdClient{},
			expectedFatal: true,
			expectedErr:   true,
		},
		{
			description: "upgrade failed should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: errors.New("boo! failed"),
			},
			oldEtcdClient: fakeTLSEtcdClient{},
			expectedFatal: true,
			expectedErr:   true,
		},
		{
			description: "wait oldEtcdClient available failed should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: errors.New("boo! failed"),
			},
			oldEtcdClient: fakeTLSEtcdClient{
				availableErr: true,
			},
			expectedFatal: true,
			expectedErr:   true,
		},
		{
			description: "wait newEtcdClient available failed should result in a rollback and an abort",
			waitErrsToReturn: map[string]error{
				waitForHashes:        nil,
				waitForHashChange:    nil,
				waitForPodsWithLabel: nil,
			},
			oldEtcdClient: fakeTLSEtcdClient{},
			newEtcdClient: fakePodManifestEtcdClient{AvailableErr: true},
			expectedFatal: true,
			expectedErr:   true,
		},
	}

	for _, rt := range tests {
		waiter := NewFakeStaticPodWaiter(rt.waitErrsToReturn)
		pathMgr, newcfg := getPathMgrAndCfg(t)
		defer pathMgr.(*fakeStaticPodPathManager).CleanupKubernetesDir()
		defer os.RemoveAll(newcfg.CertificatesDir)
		defer os.RemoveAll(newcfg.Etcd.Local.DataDir)

		recoverManifests := map[string]string{}
		if rt.newEtcdClient == nil {
			rt.newEtcdClient = fakePodManifestEtcdClient{
				ManifestDir:     pathMgr.RealManifestDir(),
				CertificatesDir: newcfg.CertificatesDir,
			}
		} else {
			rt.newEtcdClient = fakePodManifestEtcdClient{
				ManifestDir:     pathMgr.RealManifestDir(),
				CertificatesDir: newcfg.CertificatesDir,
				AvailableErr:    true,
			}
		}
		actualFatal, actualErr := performEtcdStaticPodUpgrade(
			nil,
			waiter,
			pathMgr,
			newcfg,
			recoverManifests,
			rt.oldEtcdClient,
			rt.newEtcdClient,
		)
		if (actualFatal != rt.expectedFatal) || ((actualErr != nil) != rt.expectedErr) {
			t.Errorf(
				"failed performEtcdStaticPodUpgrade\n%s\n\texpected error: %t\n\tgot: %t\n\tactual error: %v",
				rt.description,
				rt.expectedErr,
				(actualErr != nil),
				actualErr,
			)
		}
	}
}
