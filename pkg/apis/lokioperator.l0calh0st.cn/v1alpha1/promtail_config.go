package v1alpha1

// this go file is promtail configure
// all configure you can find in loki.doc website: https://grafana.com/docs/loki/latest/clients/promtail/configuration/

import "time"


type PromtailConfig struct {
	Server PromtailServerConfig `json:"server"`
	Clients PromtailClientConfig `json:"clients"`
	Positions PromtailPositionsConfig `json:"positions"`
	ScrapeConfigs ScrapeConfigs `json:"scrape_configs"`
}

// PromtailServerConfig configures Promtail's behaviour as an HTTP server
type PromtailServerConfig struct {
	Disable bool `json:"disable"`
	HttpListenAddress string `json:"http_listen_address"`
}


// PromtailClientConfig configures how promtail connects to an instance of Loki
type PromtailClientConfig struct {
	URL string `json:"url"`
	TenantId string `json:"tenant_id"`
}

// PromtailPositionsConfig configures where promtail save a file indicationg how fat it has read into a file
type PromtailPositionsConfig struct {
	Filename string `json:"filename"`		// default = "/var/log/positions.yaml
	SyncPeriod time.Duration `json:"sync_period"`
	IgnoreInvalidYaml bool `json:"ignore_invalid_yaml"`
}


// ScrapeConfigs configures how promtail can scrape logs from a series of targets using a specified discovery method
type ScrapeConfigs struct {
	// get config from scrape
}



// TargetConfig controls the behaviour of reading files from discovered targets
type TargetConfig struct {
	// Period to resync directories being watched and files being tailed to discover new ons or stop watching removed ones
	SyncPeriod time.Duration
}



