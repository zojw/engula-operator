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

	clusterv1alpha1 "github.com/engula/engula-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ClusterServiceAccountName(clusterName string) string {
	return clusterName + "-sa"
}

func ClusterRoleName(clusterName string) string {
	return clusterName + "-role"
}

func ClusterBindingName(clusterName string) string {
	return clusterName + "-binding"
}

func CheckClusterServiceAccountExist(clientset *kubernetes.Clientset, cluster *clusterv1alpha1.Cluster) (ready bool, err error) {
	ready = true
	sa := clientset.CoreV1().ServiceAccounts(cluster.Namespace)
	if _, err = sa.Get(context.Background(), ClusterServiceAccountName(cluster.Name), metav1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			ready, err = false, nil
			return
		}
	}
	return
}

var (
	_ ResourceBuilder = &ServiceAccountBuilder{}
	_ ResourceBuilder = &RoleBuilder{}
	_ ResourceBuilder = &RoleBindingBuilder{}
)

func needDeployComponent(cluster *clusterv1alpha1.Cluster) bool {
	spec := cluster.Spec
	return spec.Kernel != nil || spec.Engine != nil || spec.Storage != nil ||
		spec.Journal != nil || spec.Background != nil
}

type ServiceAccountBuilder struct {
	Cluster *clusterv1alpha1.Cluster
}

func (s *ServiceAccountBuilder) NeedBuild() bool {
	return needDeployComponent(s.Cluster)
}

func (s *ServiceAccountBuilder) QueryObject() client.Object {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ClusterServiceAccountName(s.Cluster.Name),
			Namespace: s.Cluster.Namespace,
		},
	}
}

func (s *ServiceAccountBuilder) Build(_ runtime.Object, object client.Object) (desired client.Object, err error) {
	return object, nil
}

type RoleBuilder struct {
	Cluster *clusterv1alpha1.Cluster
}

func (s *RoleBuilder) NeedBuild() bool {
	return needDeployComponent(s.Cluster)
}

func (s *RoleBuilder) QueryObject() client.Object {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ClusterRoleName(s.Cluster.Name),
			Namespace: s.Cluster.Namespace,
		},
		Rules: []rbacv1.PolicyRule{},
	}
}

func (s *RoleBuilder) Build(_ runtime.Object, object client.Object) (desired client.Object, err error) {
	return object, nil
}

type RoleBindingBuilder struct {
	Cluster *clusterv1alpha1.Cluster
}

func (s *RoleBindingBuilder) NeedBuild() bool {
	return needDeployComponent(s.Cluster)
}

func (s *RoleBindingBuilder) QueryObject() client.Object {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ClusterBindingName(s.Cluster.Name),
			Namespace: s.Cluster.Namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     ClusterRoleName(s.Cluster.Name),
		},
	}
}

func (s *RoleBindingBuilder) Build(_ runtime.Object, object client.Object) (desired client.Object, err error) {
	b := object.(*rbacv1.RoleBinding)

	clusterSA := ClusterServiceAccountName(s.Cluster.Name)
	var exist bool
	for _, s := range b.Subjects {
		if s.Name == clusterSA {
			exist = true
			break
		}
	}

	if !exist {
		b.Subjects = append(b.Subjects, rbacv1.Subject{
			Kind: "ServiceAccount",
			Name: clusterSA,
		})
	}

	return b, nil
}
