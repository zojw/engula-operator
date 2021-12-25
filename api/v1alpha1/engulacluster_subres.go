/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import corev1 "k8s.io/api/core/v1"

type SubRes struct {
	Replicas  int32
	Image     *PodImage
	Port      *int32
	Resources corev1.ResourceRequirements
	Volume    Volume
}

func (s *EngulaClusterSpec) ToSubRes() (res SubRes) {
	if s.Engine != nil {

	}
	res = append(res, s.Engine, s.Kernel, s.Storage, s.Journal, s.Background)
	return
}

func (in *Engine) Engineinto() SubRes {
	return SubRes{}
}

func (in *Kernel) into() SubRes {
	//TODO implement me
	panic("implement me")
}

func (in *Storage) into() SubRes {
	//TODO implement me
	panic("implement me")
}

func (in *Journal) into() SubRes {
	//TODO implement me
	panic("implement me")
}

func (in *Background) into() SubRes {
	//TODO implement me
	panic("implement me")
}
