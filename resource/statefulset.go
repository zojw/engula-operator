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

package resource

import (
	"github.com/engula/engula-operator/api/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

const (
	containerName = "unit"
	volumeName    = "data"
	MountPath     = "/data"
)

var _ ResourceBuilder = &StatefulsetBuilder{}

type StatefulResource struct {
	Replicas  int32
	Image     *v1alpha1.PodImage
	Port      *int32
	Resources corev1.ResourceRequirements
	Volume    v1alpha1.VolumeClaim
}

type StatefulsetBuilder struct {
	Own         func() *StatefulResource
	OwnName     string
	ClusterName string
	Component   string
	Args        []string
	Selector    Labels
}

func NewStatefulsetBuilder(ownName, clusterName string, selector Labels, f func() *StatefulResource, component Component, args ...string) *StatefulsetBuilder {
	var memo *StatefulResource
	return &StatefulsetBuilder{
		Own: func() *StatefulResource {
			if memo != nil {
				return memo
			}
			memo = f()
			return memo
		},
		Component:   string(component),
		Args:        args,
		OwnName:     ownName,
		ClusterName: clusterName,
		Selector:    selector,
	}
}

func (b *StatefulsetBuilder) NeedBuild() bool {
	return b.Own != nil
}

func (b *StatefulsetBuilder) QueryObject() client.Object {
	return &v1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: b.statefulsetName(),
		},
	}
}

func (b *StatefulsetBuilder) Build(current runtime.Object, object client.Object) (desired client.Object, err error) {
	sts := object.(*v1.StatefulSet)
	if sts.ObjectMeta.Name == "" {
		sts.ObjectMeta.Name = b.statefulsetName()
	}

	sts.Spec = v1.StatefulSetSpec{
		Replicas: &b.Own().Replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: b.Selector,
		},
		Template:            b.buildPod(),
		ServiceName:         b.headlessServiceName(),
		PodManagementPolicy: v1.ParallelPodManagement,
		UpdateStrategy: v1.StatefulSetUpdateStrategy{
			RollingUpdate: &v1.RollingUpdateStatefulSetStrategy{},
		},
	}
	if err = b.Own().Volume.Apply(volumeName, containerName, MountPath, &sts.Spec, b.Selector); err != nil {
		return
	}

	if sts.ObjectMeta.Labels == nil {
		sts.ObjectMeta.Labels = map[string]string{}
	}
	if sts.ObjectMeta.Annotations == nil {
		sts.ObjectMeta.Annotations = map[string]string{}
	}
	err = buildDesiredLabelsAnnotations(current, sts, b.Selector)
	if err != nil {
		return
	}

	desired = sts
	return
}

func (b *StatefulsetBuilder) buildPod() corev1.PodTemplateSpec {
	intPtr := func(i int64) *int64 {
		return &i
	}
	boolPtr := func(i bool) *bool {
		return &i
	}
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: b.Selector,
		},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: intPtr(60),
			Containers:                    b.buildContainers(),
			AutomountServiceAccountToken:  boolPtr(false),
			ServiceAccountName:            b.accountName(),
		},
	}
}

func (b *StatefulsetBuilder) buildContainers() []corev1.Container {
	image := b.Own().Image
	return []corev1.Container{{
		Name:            containerName,
		Image:           image.Name,
		ImagePullPolicy: *image.PullPolicyName,
		Resources:       b.Own().Resources,
		Command:         b.buildCmd(),
		Env:             b.buildEnv(),
		Ports: []corev1.ContainerPort{
			{
				Name:          containerName,
				ContainerPort: *b.Own().Port,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		ReadinessProbe: nil, // TODO!!!!
		Lifecycle:      nil, // TODO!!!!
	}}
}

func (b *StatefulsetBuilder) buildEnv() []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name: "POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name:  "kubecluster",
			Value: b.ClusterName,
		},
		{
			Name:  "journalport",
			Value: strconv.Itoa(int(v1alpha1.DefaultJournalPort)), //TODO:...it's not a good way
		},
		{
			Name:  "storageport",
			Value: strconv.Itoa(int(v1alpha1.DefaultStoragePort)), //TODO:...it's not a good way
		},
		{
			Name:  "kernelport",
			Value: strconv.Itoa(int(v1alpha1.DefaultKernelPort)), //TODO:...it's not a good way
		},
	}
}

func (b *StatefulsetBuilder) buildCmd() []string {
	cmd := []string{
		"/bin/engula",
		b.Component,
		"run",
	}
	cmd = append(cmd, b.Args...)
	return cmd
}

func (b *StatefulsetBuilder) accountName() string {
	return ClusterServiceAccountName(b.ClusterName)
}

func (b *StatefulsetBuilder) statefulsetName() string {
	return b.OwnName
}

func (b *StatefulsetBuilder) headlessServiceName() string {
	return b.OwnName
}
