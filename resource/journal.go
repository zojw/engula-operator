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

var _ ResourceBuilder = &JournalBuilder{}

type JournalBuilder struct {
	Own      *v1alpha1.Cluster
	Selector Labels
}

func (b *JournalBuilder) NeedBuild() bool {
	return b.Own.Spec.Journal != nil
}

func (b *JournalBuilder) QueryObject() client.Object {
	return JournalQueryObject(b.name())
}

func JournalQueryObject(name string) *v1alpha1.Journal {
	return &v1alpha1.Journal{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func (b *JournalBuilder) name() string {
	return b.Own.Name + "-journal"
}

func (b *JournalBuilder) Build(current runtime.Object, object client.Object) (desired client.Object, err error) {
	e := object.(*v1alpha1.Journal)
	if e.ObjectMeta.Name == "" {
		e.ObjectMeta.Name = b.name()
	}

	e.Spec = *b.Own.Spec.Journal.DeepCopy()
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
