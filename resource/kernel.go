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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ ResourceBuilder = &KernelBuilder{}

type KernelBuilder struct {
	Own      *v1alpha1.Cluster
	Selector Labels
}

func (b *KernelBuilder) NeedBuild() bool {
	return b.Own.Spec.Kernel != nil
}

func (b *KernelBuilder) QueryObject() client.Object {
	return KernelQueryObject(b.name())
}

func KernelQueryObject(name string) *v1alpha1.Kernel {
	return &v1alpha1.Kernel{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func (b *KernelBuilder) name() string {
	return b.Own.Name + "-kernel"
}

func (b *KernelBuilder) Build(current runtime.Object, object client.Object) (desired client.Object, err error) {
	k := object.(*v1alpha1.Kernel)
	if k.ObjectMeta.Name == "" {
		k.ObjectMeta.Name = b.name()
	}

	k.Spec = *b.Own.Spec.Kernel.DeepCopy()
	if k.Spec.Image == nil && b.Own.Spec.Image != nil {
		k.Spec.Image = b.Own.Spec.Image.DeepCopy()
	}

	if k.ObjectMeta.Labels == nil {
		k.ObjectMeta.Labels = map[string]string{}
	}
	if k.ObjectMeta.Annotations == nil {
		k.ObjectMeta.Annotations = map[string]string{}
	}
	err = buildDesiredLabelsAnnotations(current, k, b.Selector)
	if err != nil {
		return
	}

	desired = k
	return
}
