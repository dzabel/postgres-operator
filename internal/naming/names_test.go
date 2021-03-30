/*
 Copyright 2021 Crunchy Data Solutions, Inc.
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

package naming

import (
	"testing"

	"gotest.tools/v3/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crunchydata/postgres-operator/pkg/apis/postgres-operator.crunchydata.com/v1alpha1"
)

func TestAsObjectKey(t *testing.T) {
	assert.Equal(t, AsObjectKey(
		metav1.ObjectMeta{Namespace: "ns1", Name: "thing"}),
		client.ObjectKey{Namespace: "ns1", Name: "thing"})
}

func TestClusterNamesUniqueAndValid(t *testing.T) {
	cluster := &v1alpha1.PostgresCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "ns1", Name: "pg0",
		},
	}

	type test struct {
		name  string
		value metav1.ObjectMeta
	}

	t.Run("ConfigMaps", func(t *testing.T) {
		names := sets.NewString()
		for _, tt := range []test{
			{"ClusterConfigMap", ClusterConfigMap(cluster)},
			{"PatroniDistributedConfiguration", PatroniDistributedConfiguration(cluster)},
			{"PatroniLeaderConfigMap", PatroniLeaderConfigMap(cluster)},
			{"PatroniTrigger", PatroniTrigger(cluster)},
		} {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.value.Namespace, cluster.Namespace)
				assert.Assert(t, tt.value.Name != cluster.Name, "may collide")
				assert.Assert(t, !names.Has(tt.value.Name), "%q defined already", tt.value.Name)
				assert.Assert(t, nil == validation.IsDNS1123Label(tt.value.Name))
				names.Insert(tt.value.Name)
			})
		}
	})

	t.Run("Secrets", func(t *testing.T) {
		names := sets.NewString()
		for _, tt := range []test{
			{"PostgresUserSecret", PostgresUserSecret(cluster)},
			{"PostgresTLSSecret", PostgresTLSSecret(cluster)},
			{"PatroniAuthSecret", PatroniAuthSecret(cluster)},
		} {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.value.Namespace, cluster.Namespace)
				assert.Assert(t, tt.value.Name != cluster.Name, "may collide")
				assert.Assert(t, !names.Has(tt.value.Name), "%q defined already", tt.value.Name)
				assert.Assert(t, nil == validation.IsDNS1123Label(tt.value.Name))
				names.Insert(tt.value.Name)
			})
		}
	})

	t.Run("Services", func(t *testing.T) {
		names := sets.NewString()
		for _, tt := range []test{
			{"ClusterPodService", ClusterPodService(cluster)},
			{"ClusterPrimaryService", ClusterPrimaryService(cluster)},
			// Patroni can use Endpoints which relate directly to a Service.
			{"PatroniDistributedConfiguration", PatroniDistributedConfiguration(cluster)},
			{"PatroniLeaderEndpoints", PatroniLeaderEndpoints(cluster)},
			{"PatroniTrigger", PatroniTrigger(cluster)},
		} {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.value.Namespace, cluster.Namespace)
				assert.Assert(t, tt.value.Name != cluster.Name, "may collide")
				assert.Assert(t, !names.Has(tt.value.Name), "%q defined already", tt.value.Name)
				assert.Assert(t, nil == validation.IsDNS1123Label(tt.value.Name))
				names.Insert(tt.value.Name)
			})
		}
	})
}

func TestInstanceNamesUniqueAndValid(t *testing.T) {
	instance := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "ns", Name: "some-such",
		},
	}

	type test struct {
		name  string
		value metav1.ObjectMeta
	}

	t.Run("ConfigMaps", func(t *testing.T) {
		names := sets.NewString()
		for _, tt := range []test{
			{"InstanceConfigMap", InstanceConfigMap(instance)},
		} {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.value.Namespace, instance.Namespace)
				assert.Assert(t, tt.value.Name != instance.Name, "may collide")
				assert.Assert(t, !names.Has(tt.value.Name), "%q defined already", tt.value.Name)
				assert.Assert(t, nil == validation.IsDNS1123Label(tt.value.Name))
				names.Insert(tt.value.Name)
			})
		}
	})

	t.Run("Secrets", func(t *testing.T) {
		names := sets.NewString()
		for _, tt := range []test{
			{"InstanceCertificates", InstanceCertificates(instance)},
		} {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.value.Namespace, instance.Namespace)
				assert.Assert(t, tt.value.Name != instance.Name, "may collide")
				assert.Assert(t, !names.Has(tt.value.Name), "%q defined already", tt.value.Name)
				assert.Assert(t, nil == validation.IsDNS1123Label(tt.value.Name))
				names.Insert(tt.value.Name)
			})
		}
	})
}
