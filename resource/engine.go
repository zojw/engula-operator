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

var _ ResourceBuilder = &EngineBuilder{}

type EngineBuilder struct {
	Own      *v1alpha1.Cluster
	Selector Labels
}

func (b *EngineBuilder) NeedBuild() bool {
	return b.Own.Spec.Engine != nil
}

func (b *EngineBuilder) QueryObject() client.Object {
	return EngineQueryObject(b.name())
}

func EngineQueryObject(name string) *v1alpha1.Engine {
	return &v1alpha1.Engine{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func (b *EngineBuilder) name() string {
	return b.Own.Name + "-engine"
}

func (b *EngineBuilder) Build(current runtime.Object, object client.Object) (desired client.Object, err error) {
	e := object.(*v1alpha1.Engine)
	if e.ObjectMeta.Name == "" {
		e.ObjectMeta.Name = b.name()
	}

	e.Spec = *b.Own.Spec.Engine.DeepCopy()
	if e.Spec.Image == nil && b.Own.Spec.Image != nil {
		e.Spec.Image = b.Own.Spec.Image.DeepCopy()
	}

	if e.ObjectMeta.Labels == nil {
		e.ObjectMeta.Labels = map[string]string{}
	}
	if e.ObjectMeta.Annotations == nil {
		e.ObjectMeta.Annotations = map[string]string{}
	}
	err = buildDesiredLabelsAnnotations(current, e, b.Selector)
	if err != nil {
		return
	}

	desired = e
	return
}
