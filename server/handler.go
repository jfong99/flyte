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

package server

import (
	"github.com/ExpediaGroup/flyte/audit"
	"github.com/ExpediaGroup/flyte/datastore"
	"github.com/ExpediaGroup/flyte/execution"
	"github.com/ExpediaGroup/flyte/flow"
	"github.com/ExpediaGroup/flyte/flytepath"
	"github.com/ExpediaGroup/flyte/httputil"
	"github.com/ExpediaGroup/flyte/info"
	"github.com/ExpediaGroup/flyte/pack"
	"github.com/husobee/vestigo"
	"net/http"
)

func Handler() http.Handler {

	router := vestigo.NewRouter()

	// --- swagger ---
	swaggerUi := http.FileServer(http.Dir("swagger/swagger-ui"))
	router.Handle("/swagger", http.StripPrefix("/swagger", swaggerUi))
	router.Handle("/swagger/:file", http.StripPrefix("/swagger", swaggerUi))
	swaggerConfig := http.FileServer(http.Dir("swagger"))
	router.Handle("/swagger-config/:file", http.StripPrefix("/swagger-config", swaggerConfig))

	// --- info ---
	router.Get(flytepath.IndexPath, info.Index)
	router.Get(flytepath.VersionPath, info.V1)
	router.Get(flytepath.HealthPath, info.Health)
	router.Get(flytepath.VersionDocPath, info.V1Swagger)

	// --- pack ---
	router.Get(flytepath.PacksPath, pack.GetPacks)
	router.Post(flytepath.PacksPath, pack.PostPack, YamlHandler)
	router.Get(flytepath.PackPath, pack.GetPack)
	router.Delete(flytepath.PackPath, pack.DeletePack)

	// --- execution ---
	router.Post(flytepath.TakeActionPath, execution.TakeAction, YamlHandler)
	router.Post(flytepath.PostEventPath, execution.PostEvent, YamlHandler)
	router.Post(flytepath.TakeActionResultPath, execution.CompleteAction, YamlHandler)

	// --- flow ---
	router.Get(flytepath.FlowsPath, flow.GetFlows)
	router.Post(flytepath.FlowsPath, flow.PostFlow, YamlHandler)
	router.Get(flytepath.FlowPath, flow.GetFlow)
	router.Delete(flytepath.FlowPath, flow.DeleteFlow)

	// --- datastore ---
	router.Get(flytepath.DatastorePath, datastore.GetItems)
	router.Get(flytepath.DatastoreItemPath, datastore.GetItem)
	router.Put(flytepath.DatastoreItemPath, datastore.StoreItem, YamlHandler)
	router.Delete(flytepath.DatastoreItemPath, datastore.DeleteItem)

	// --- audit ---
	router.Get(flytepath.AuditFlowPath, audit.GetFlows)
	router.Get(flytepath.AuditGetFlow, audit.GetFlow)

	return wrapRequestInterceptorAround(router)
}

func wrapRequestInterceptorAround(h http.Handler) http.Handler {
	return httputil.NewRequestInterceptor(h)
}
