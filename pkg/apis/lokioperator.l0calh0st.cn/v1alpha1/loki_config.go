package v1alpha1

import "time"

// LokiServerConfig is used to configures the server of the launched module(s).
// It is referenced to server in loki yaml config
type LokiServerConfig struct {
	HttpListenAddress              string        `json:"http_listen_address"`
	HttpListenPort                 int           `json:"http_listen_port"`
	GrpcListenAddress              string        `json:"grpc_listen_address"`
	GrpcListenPort                 int           `json:"grpc_listen_port"`
	RegisterInstrumentation        bool          `json:"register_instrumentation"`
	GracefulShutdownTimeout        time.Duration `json:"graceful_shutdown_timeout"`
	HtpServerReadTimeout           time.Duration `json:"htp_server_read_timeout"`
	HttpServerWriteTimeout         time.Duration `json:"http_server_write_timeout"`
	HttpServerIdleTimeout          time.Duration `json:"http_server_idle_timeout"`
	GrpcServerMaxRecvMsgSize       int           `json:"grpc_server_max_recv_msg_size"`
	GrpcServerMaxSendMsgSize       int           `json:"grpc_server_max_send_msg_size"`
	GrpcServerMaxConcurrentStreams int           `json:"grpc_server_max_concurrent_streams"`
	LogLevel                       string        `json:"log_level"`
	HttpPathPrefix                 string        `json:"http_path_prefix"`
}

// DistributorConfig is used to configure the distributor.
// It is references to distributor
type DistributorConfig struct {
}

// QuerierConfig is used to configures the querier
type QuerierConfig struct {
}

// QueryScheduler is used to the Loki query scheduler
type QuerySchedulerConfig struct {
}

// QueryFrontendConfig is used to configures the Loki query-frontend
type QueryFrontendConfig struct {
}

// QueryrangeConfig is used to configures the  the query splitting and caching in the Loki query-frontend
type QueryrangeConfig struct {
}

// RulerConfig is used to config the Loki ruler
type RulerConfig struct {
}

// IngesterClientConfig is used to configures how the distributor will connect to ingesters
// Only appropriate when running all modules, the distributor, or the querier.
type IngesterClientConfig struct {
}

// IngesterConfig is used to configures ingester and how the ingester will register itself to a kv store
type IngesterConfig struct {
}

// StorageConfig is used to configures where Loki will store data
type StorageConfig struct {
}

// ChunkStoreConfig is used to how Loki will store data in the specific store.
type ChunkStoreConfig struct {
}

// SchemaConfig configures the chunk index schema and where it is stored.
type SchemaConfig struct {
}

//CompactorConfig configures the compactor component which compacts index shards for performance.
type CompactorConfig struct {
}

// LimitsConfig configures limits per-tenant or globally.
type LimitsConfig struct {
}

// The frontend_worker_config configures the worker - running within the Loki
//# querier - picking up and executing queries enqueued by the query-frontend.
type FrontendWorkerConfig struct {
}

// TableManagerConfig configures the table manager for retention.
type TableManagerConfig struct {
}

// RuntimeConfig  for "runtime config" module, responsible for reloading runtime configuration file.
type RuntimeConfig struct {
	// File is used to periodiccally check and reload
	File   string        `json:"file"`
	Period time.Duration `json:"period"`
}

// TracingConfig is used to config tracing
type TracingConfig struct {
}

// Common config to be shared between multiple modules
type CommonConfig struct {
	PathPrefix        string `json:"path_prefix,omitempty"`
	ReplicationFactor int    `json:"replication_factor,omitempty"`
	PersistTokens     bool   `json:"persist_tokens,omitempty"`
}
