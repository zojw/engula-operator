// Copyright 2021 The Engula Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"errors"

	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (v *VolumeClaim) Apply(volumeName, containerName, mountPath string, sts *v1.StatefulSetSpec, selector map[string]string) error {
	if v == nil {
		return nil
	}
	var c *corev1.Container
	for i := range sts.Template.Spec.Containers {
		tmp := &sts.Template.Spec.Containers[i]
		if tmp.Name == containerName {
			c = tmp
			break
		}
	}
	if c == nil {
		return errors.New("container: " + containerName + "not found")
	}
	c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
		Name:      volumeName,
		MountPath: mountPath,
	})

	volume := corev1.Volume{
		Name: volumeName,
	}
	volume.VolumeSource = corev1.VolumeSource{
		PersistentVolumeClaim: &v.PersistentVolumeSource,
	}
	sts.Template.Spec.Volumes = append(sts.Template.Spec.Volumes, volume)

	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:   volumeName,
			Labels: selector,
		},
		Spec: v.PersistentVolumeClaimSpec,
	}
	sts.VolumeClaimTemplates = append(sts.VolumeClaimTemplates, pvc)

	return nil
}
