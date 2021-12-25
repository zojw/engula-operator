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
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const LastAppliedAnnotation = "engine.io/last-applied"

var (
	lastAppliedAnnotator  = patch.NewAnnotator(LastAppliedAnnotation)
	lastAppliedPatchMaker = patch.NewPatchMaker(lastAppliedAnnotator)
)

type Resource struct {
	client.Client

	ctx context.Context
	ns  string
}

func NewResource(ctx context.Context, ns string, c client.Client) *Resource {
	return &Resource{
		Client: c,
		ctx:    ctx,
		ns:     ns,
	}
}

func (r *Resource) Fetch(o client.Object) error {
	accessor, err := meta.Accessor(o)
	if err != nil {
		return err
	}
	key := types.NamespacedName{
		Name:      accessor.GetName(),
		Namespace: r.ns,
	}
	return r.Get(r.ctx, key, o)
}

func (r *Resource) Apply(obj client.Object, fn func(client.Object) (client.Object, error)) (updated bool, err error) {
	var (
		desired  client.Object
		changed  bool
		accessor metav1.Object
	)

	accessor, err = meta.Accessor(obj)
	if err != nil {
		return
	}
	accessor.SetNamespace(r.ns)

	key := client.ObjectKeyFromObject(obj)
	if err = r.Get(r.ctx, key, obj); err != nil {
		if !errors.IsNotFound(err) {
			return
		}
		desired, err = fn(obj)
		if err != nil {
			return
		}
		err = lastAppliedAnnotator.SetLastAppliedAnnotation(desired)
		if err != nil {
			return
		}
		if err = r.Create(r.ctx, desired); err != nil {
			return
		}
		updated = true
		return
	}

	origin := obj.DeepCopyObject()
	desired, err = fn(obj)
	if err != nil {
		return
	}
	changed, err = IsChanged(origin, desired)
	if err != nil || changed {
		return
	}
	err = lastAppliedAnnotator.SetLastAppliedAnnotation(desired)
	if err != nil {
		return
	}
	if err = r.Update(r.ctx, desired); err != nil {
		return
	}
	updated = true
	return
}

func IsChanged(origin, desired runtime.Object) (bool, error) {
	opts := []patch.CalculateOption{
		patch.IgnoreStatusFields(),
	}

	switch desired.(type) {
	case *appsv1.StatefulSet:
		opts = append(opts, patch.IgnoreVolumeClaimTemplateTypeMetaAndStatus())
	}

	patchResult, err := lastAppliedPatchMaker.Calculate(origin, desired, opts...)
	if err != nil {
		return false, err
	}
	return !patchResult.IsEmpty(), nil
}
