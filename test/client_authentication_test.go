// Copyright (c) 2008-2018, Hazelcast, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License")
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

package test

import (
	"testing"

	hazelcast "github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/security"
	"github.com/hazelcast/hazelcast-go-client/serialization"
	"github.com/hazelcast/hazelcast-go-client/test/assert"
)

var samplePortableFactoryID int32 = 666

func TestCustomAuthentication(t *testing.T) {
	cluster, _ := remoteController.CreateCluster("", DefaultServerConfig)
	remoteController.StartMember(cluster.ID)
	defer remoteController.ShutdownCluster(cluster.ID)

	cfg := hazelcast.NewConfig()
	cfg.SerializationConfig().AddPortableFactory(samplePortableFactoryID, &portableFactory{})

	cfg.SecurityConfig().SetCredentials(&CustomCredentials{
		security.NewUsernamePasswordCredentials(
			"dev",
			"dev-pass",
		),
	})

	client, _ := hazelcast.NewClientWithConfig(cfg)
	defer client.Shutdown()
	mp, _ := client.GetMap("myMap")
	_, err := mp.Put("key", "value")

	assert.ErrorNil(t, err)
}

type portableFactory struct {
}

func (pf *portableFactory) Create(classID int32) serialization.Portable {
	if classID == samplePortableFactoryID {
		return &CustomCredentials{}
	}
	return nil
}

type CustomCredentials struct {
	*security.UsernamePasswordCredentials
}

func (cc *CustomCredentials) FactoryID() (factoryID int32) {
	return samplePortableFactoryID
}

func (cc *CustomCredentials) ClassID() (classID int32) {
	return 7
}