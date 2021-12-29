/*
Copyright 2021 The loki-operator Authors.
Licensed under the Apache License, PROJECT_VERSION 2.0 (the "License");
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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ModeKind string

const (
	ModeKindMonolithic     ModeKind = "monolithic"
	ModeKindSampleScalable ModeKind = "sample"
	ModeKinMicroservice    ModeKind = "microservice"
)

type LokiTargetKind string

const (
	TargetKindAllInOne         LokiTargetKind = "all"
	TargetKindReadAsRead       LokiTargetKind = "read"
	TargetKindAsWrite          LokiTargetKind = "write"
	TargetKindAsIngester       LokiTargetKind = "ingester"
	TargetKindAsDistributor    LokiTargetKind = "distributor"
	TargetKindAsQueryFrontend  LokiTargetKind = "query-frontent"
	TargetKindAsQueryScheduler LokiTargetKind = "query-scheduler"
	TargetKindAsQuerier        LokiTargetKind = "querier"
	TargetKindAsIndexGateway   LokiTargetKind = "index-gateway"
	TargetKindAsRuler          LokiTargetKind = "ruler"
	TargetKindAsCompactor      LokiTargetKind = "compactor"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true

// Loki defines Loki deployment
type Loki struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LokiSpec   `json:"spec"`
	Status LokiStatus `json:"status"`
}

// LokiSpec describes the specification of Loki applications using kubernetes as a cluster manager
type LokiSpec struct {
	Image      string                                `json:"image"`
	DeployMode map[ModeKind]map[LokiTargetKind]int32 `json:"deployMode,omitempty"`
	AuthEnable bool                                  `json:"authEnable"`
	ConfigMap  string                                `json:"configMap,omitempty"`
	Config     struct {
		Server           LokiServerConfig     `json:"server,omitempty"`
		Distributor      DistributorConfig    `json:"distributor,omitempty"`
		Querier          QuerierConfig        `json:"querier,omitempty"`
		QueryScheduler   QuerySchedulerConfig `json:"query_scheduler,omitempty"`
		Frontend         QueryFrontendConfig  `json:"frontend,omitempty"`
		QueryRange       QueryrangeConfig     `json:"query_range,omitempty"`
		Ruler            RulerConfig          `json:"ruler,omitempty"`
		IngesterClient   IngesterClientConfig `json:"ingester_client,omitempty"`
		Ingester         IngesterConfig       `json:"ingester,omitempty"`
		StorageConfig    StorageConfig        `json:"storage_config,omitempty"`
		ChunkStoreConfig ChunkStoreConfig     `json:"chunk_store_config,omitempty"`
		SchemaConfig     SchemaConfig         `json:"schema_config,omitempty"`
		Compactor        CompactorConfig      `json:"compactor,omitempty"`
		LimitsConfig     LimitsConfig         `json:"limits_config,omitempty"`
		FrontendWorker   FrontendWorkerConfig `json:"frontend_worker,omitempty"`
		TableManager     TableManagerConfig   `json:"table_manager,omitempty"`
		RuntimeConfig    RuntimeConfig        `json:"runtime_config,omitempty"`
		Tracing          TracingConfig        `json:"tracing,omitempty"`
		Common           CommonConfig         `json:"common,omitempty"`
	} `json:"config"`
}

// LokiStatus describes the current status of Loki applications
type LokiStatus struct {
	// todo, write your code
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LokiList carries a list of Loki objects
type LokiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Loki `json:"items"`
}
