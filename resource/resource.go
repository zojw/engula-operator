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
	"github.com/banzaicloud/k8s-objectmatcher/patch"
	clusterv1alpha1 "github.com/engula/engula-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const LastAppliedAnnotation = "engine.io/last-applied"

var (
	lastAppliedAnnotator  = patch.NewAnnotator(LastAppliedAnnotation)
	lastAppliedPatchMaker = patch.NewPatchMaker(lastAppliedAnnotator)
)

type Labels map[string]string

type ManagedResources struct {
	client.Client
	Selector Labels

	ctx context.Context
	ns  string
}

func NewManagedResources(ctx context.Context, ns string, c client.Client) *ManagedResources {
	return &ManagedResources{
		Client: c,
		ctx:    ctx,
		ns:     ns,
	}
}

func ClusterQueryObject(name string) *clusterv1alpha1.Cluster {
	return &clusterv1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func (r *ManagedResources) Fetch(o client.Object) error {
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

func (r *ManagedResources) apply(desired client.Object, needCreate bool) (err error) {
	if err = addNs(desired, r.ns); err != nil {
		return
	}
	if needCreate {
		if err = lastAppliedAnnotator.SetLastAppliedAnnotation(desired); err != nil {
			return
		}
		err = r.Create(r.ctx, desired)
		return
	}
	if err = lastAppliedAnnotator.SetLastAppliedAnnotation(desired); err != nil {
		return
	}
	err = r.Update(r.ctx, desired)
	return
}

func addNs(o runtime.Object, ns string) error {
	accessor, err := meta.Accessor(o)
	if err != nil {
		return err
	}
	accessor.SetNamespace(ns)
	return nil
}

func buildDesiredLabelsAnnotations(current, desired runtime.Object, internal Labels) (err error) {
	var (
		currentAccessor metav1.Object
		desiredAccessor metav1.Object
	)
	currentAccessor, err = meta.Accessor(current)
	if err != nil {
		return
	}
	desiredAccessor, err = meta.Accessor(desired)
	if err != nil {
		return
	}
	labels := make(map[string]string)
	for k, v := range currentAccessor.GetLabels() {
		labels[k] = v
	}
	for k, v := range internal {
		labels[k] = v
	}
	for k, v := range desiredAccessor.GetLabels() {
		labels[k] = v
	}

	desiredAccessor.SetLabels(labels)

	if currentAccessor.GetAnnotations() == nil {
		return
	}
	desiredAnnotations := desiredAccessor.GetAnnotations()
	if desiredAnnotations == nil {
		desiredAnnotations = make(map[string]string)
		desiredAccessor.SetAnnotations(desiredAnnotations)
	}
	for k, v := range currentAccessor.GetAnnotations() {
		if _, exist := desiredAnnotations[k]; !exist {
			desiredAnnotations[k] = v
		}
	}
	return
}

func (r *ManagedResources) ReconcileOne(b ResourceBuilder) (updated bool, err error) {
	if !b.NeedBuild() {
		return
	}
	var needCreate bool
	query := b.QueryObject()
	if err = r.Fetch(query); err != nil {
		if !errors.IsNotFound(err) {
			return
		}
		needCreate = true
	}

	current := query.DeepCopyObject()
	var desired client.Object
	if desired, err = b.Build(current, query); err != nil {
		return
	}

	var changed bool
	if changed, err = isChanged(desired, current, needCreate); err != nil || !changed {
		return
	}

	if err = r.apply(desired, needCreate); err != nil {
		return
	}
	updated = true
	return
}

func isChanged(origin, desired runtime.Object, needCreate bool) (bool, error) {
	if needCreate {
		return true, nil
	}
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

type ResourceBuilder interface {
	NeedBuild() bool
	QueryObject() client.Object
	Build(current runtime.Object, object client.Object) (desired client.Object, err error)
}
