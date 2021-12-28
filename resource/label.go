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
	"errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var labels map[string]string

// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
const (
	NameKey      = "app.kubernetes.io/name"
	InstanceKey  = "app.kubernetes.io/instance"
	ComponentKey = "app.kubernetes.io/component"
	ManagedByKey = "app.kubernetes.io/managed-by"
)

type Component string

const (
	Kernel     Component = "kernel"
	Engine               = "engine"
	Journal              = "journal"
	Storage              = "storage"
	Background           = "background"
	Cluster              = "cluster"
)

func MakeLabels(instance string, component Component) map[string]string {
	labels = make(map[string]string)
	labels[NameKey] = "engula"
	labels[InstanceKey] = instance
	labels[ComponentKey] = string(component)
	labels[ManagedByKey] = "engula-operator"
	return labels
}

func GetClusterName(meta v1.ObjectMeta) (n string, err error) {
	cluster, exist := meta.Labels[InstanceKey]
	if !exist {
		err = errors.New("cluster name not found")
		return
	}
	n = cluster
	return
}
