/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package flexvolume

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/util/mount"
	"k8s.io/kubernetes/pkg/volume"
)

func getSecrets(spec *volume.Spec, host volume.VolumeHost) (map[string]string, error) {
	fv, _ := getVolumeSource(spec)
	secrets := make(map[string]string)
	if fv.SecretRef != nil {
		kubeClient := host.GetKubeClient()
		if kubeClient == nil {
			return nil, fmt.Errorf("Cannot get kube client")
		}

		secret, err := kubeClient.Core().Secrets(spec.Namespace).Get(fv.SecretRef.Name)
		if err != nil {
			err = fmt.Errorf("Couldn't get secret %v/%v err: %v", spec.Namespace, fv.SecretRef.Name, err)
			return nil, err
		}
		for name, data := range secret.Data {
			secrets[name] = base64.StdEncoding.EncodeToString(data)
			glog.V(1).Infof("found flex volume secret info: %s", name)
		}
	}
	return secrets, nil
}

func getVolumeSource(spec *volume.Spec) (volumeSource *api.FlexVolumeSource, readOnly bool) {
	if spec.Volume != nil && spec.Volume.FlexVolume != nil {
		volumeSource = spec.Volume.FlexVolume
		readOnly = volumeSource.ReadOnly
	} else if spec.PersistentVolume != nil {
		volumeSource = spec.PersistentVolume.Spec.FlexVolume
		readOnly = spec.ReadOnly
	}
	return
}

func prepareForMount(mounter mount.Interface, deviceMountPath string) (alreadyMounted bool, err error) {
	if _, err := os.Stat(deviceMountPath); os.IsNotExist(err) {
		if err := os.MkdirAll(deviceMountPath, 0750); err != nil {
			return false, err
		}
		return false, nil
	} else if err != nil {
		glog.Errorf("Failed to stat %s: %v", deviceMountPath, err)
		return false, err
	}
	notMnt, err := isNotMounted(mounter, deviceMountPath)
	if err != nil {
		return false, err
	}
	return !notMnt, nil
}

// Mounts the device at the given path.
// It is expected that prepareForMount has been called before.
func doMount(mounter mount.Interface, devicePath, deviceMountPath, fsType string, options []string) error {
	err := mounter.Mount(devicePath, deviceMountPath, fsType, options)
	if err != nil {
		glog.Errorf("Failed to mount the volume at %s, device: %s, error: %s", deviceMountPath, devicePath, err.Error())
		return err
	}
	return nil
}

func isNotMounted(mounter mount.Interface, deviceMountPath string) (bool, error) {
	notmnt, err := mounter.IsLikelyNotMountPoint(deviceMountPath)
	if err != nil {
		glog.Errorf("Error checking mount point %s, error: %v", deviceMountPath, err)
		return false, err
	}
	return notmnt, nil
}

// Unmount from the given path.
// The caller is expected to have checked it is mounted before.
func doUnmount(mounter mount.Interface, deviceMountPath string) error {
	if err := mounter.Unmount(deviceMountPath); err != nil {
		glog.Errorf("Failed to unmount volume: %s, error: %s", deviceMountPath, err.Error())
		return err
	}
	return nil
}

func removeMountPath(mounter mount.Interface, deviceMountPath string) error {
	notmnt, err := mounter.IsLikelyNotMountPoint(deviceMountPath)
	if err != nil {
		glog.Errorf("Error checking mount point %s, error: %v", deviceMountPath, err)
		return err
	}
	if notmnt {
		err := os.Remove(deviceMountPath)
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}
