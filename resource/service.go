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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ ResourceBuilder = &ServiceBuilder{}

type ServiceResource struct {
	Port *int32
}

type ServiceBuilder struct {
	Own      func() *ServiceResource
	OwnName  string
	Selector Labels
}

func NewServiceBuilder(ownName string, selector Labels, f func() *ServiceResource) *ServiceBuilder {
	var memo *ServiceResource
	return &ServiceBuilder{
		Own: func() *ServiceResource {
			if memo != nil {
				return memo
			}
			memo = f()
			return memo
		},
		OwnName:  ownName,
		Selector: selector,
	}
}

func (b *ServiceBuilder) NeedBuild() bool {
	return b.Own != nil
}

func (b *ServiceBuilder) QueryObject() client.Object {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: b.headlessServiceName(),
		},
	}
}

func (b *ServiceBuilder) Build(current runtime.Object, object client.Object) (desired client.Object, err error) {
	s := object.(*corev1.Service)
	if s.ObjectMeta.Name == "" {
		s.ObjectMeta.Name = b.headlessServiceName()
	}
	if s.ObjectMeta.Labels == nil {
		s.ObjectMeta.Labels = map[string]string{}
	}
	if s.ObjectMeta.Annotations == nil {
		s.ObjectMeta.Annotations = map[string]string{}
	}

	own := b.Own()
	s.Spec = corev1.ServiceSpec{
		ClusterIP:                "None",
		PublishNotReadyAddresses: true,
		Ports: []corev1.ServicePort{
			{Name: "service", Port: *own.Port},
		},
		Selector: b.Selector,
	}

	err = buildDesiredLabelsAnnotations(current, s, b.Selector)
	if err != nil {
		return
	}

	desired = s
	return
}

func (b *ServiceBuilder) headlessServiceName() string {
	return b.OwnName
}
