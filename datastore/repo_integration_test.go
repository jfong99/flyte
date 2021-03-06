// +build integration

/*
Copyright (C) 2018 Expedia Group.

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

package datastore

import (
	"github.com/ExpediaGroup/flyte/httputil"
	"github.com/ExpediaGroup/flyte/mongo"
	"github.com/ExpediaGroup/flyte/mongo/mongotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
	"os"
	"testing"
)

var mongoT *mongotest.MongoT

func TestMain(m *testing.M) {
	os.Exit(runTestsWithMongo(m))
}

func runTestsWithMongo(m *testing.M) int {
	mongoT = mongotest.NewMongoT(mongo.DbName)
	defer mongoT.Teardown()

	mongoT.Start()

	mongo.InitSession(mongoT.GetUrl(), 0)

	return m.Run()
}

func TestStore_ShouldAddNewItem(t *testing.T) {
	mongoT.DropDatabase(t)

	expected := DataItem{
		Key:         "new-item",
		Description: "My shiny new item",
		ContentType: httputil.MediaTypeJson,
		Value:       []byte(`"hello"`),
	}
	updated, err := datastoreRepo.Store(expected)
	require.NoError(t, err)

	assert.False(t, updated)
	assert.Equal(t, 1, mongoT.Count(t, mongo.DatastoreCollectionId))
	assert.Equal(t, expected, findDataItem(t, "new-item"))
}

func TestStore_ShouldUpdateExistingItem(t *testing.T) {
	mongoT.DropDatabase(t)
	existingItem := DataItem{
		Key:         "existing-item",
		Description: "My existing item",
		ContentType: httputil.MediaTypeJson,
		Value:       []byte(`"hello"`),
	}
	mongoT.Insert(t, mongo.DatastoreCollectionId, existingItem)

	assert.Equal(t, 1, mongoT.Count(t, mongo.DatastoreCollectionId))
	assert.Equal(t, existingItem, findDataItem(t, "existing-item"))

	updatedItem := DataItem{
		Key:         existingItem.Key,
		Description: "",
		ContentType: httputil.MediaTypeYaml,
		Value:       []byte(`goodbye`),
	}
	updated, err := datastoreRepo.Store(updatedItem)
	require.NoError(t, err)

	assert.True(t, updated)
	assert.Equal(t, 1, mongoT.Count(t, mongo.DatastoreCollectionId))
	assert.Equal(t, updatedItem, findDataItem(t, "existing-item"))
}

func findDataItem(t *testing.T, key string) DataItem {
	var d DataItem
	mongoT.FindOneT(t, mongo.DatastoreCollectionId, bson.M{"_id": key}, &d)
	return d
}
