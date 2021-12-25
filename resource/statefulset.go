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
	"github.com/engula/engula-operator/controllers"
)

var _ controllers.ResourceReconciler = &StatefulsetReconclier{}

type StatefulsetReconclier struct {
	statefulsetName string
	cr              *v1alpha1.EngulaCluster
}

func (s *StatefulsetReconclier) Reconcile() (updated bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *StatefulsetReconclier) IsChanged() (bool, error) {
	//TODO implement me
	panic("implement me")
}
