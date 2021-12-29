package loki

import (
	"bytes"
	crapiv1alpha1 "github.com/l0calh0st/loki-operator/pkg/apis/lokioperator.l0calh0st.cn/v1alpha1"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"text/template"
)

const lokiSimpleConfigTplForMonolithic = `auth_enabled: false
chunk_store_config:
  max_look_back_period: 0s
compactor:
  shared_store: filesystem
  working_directory: /data/loki/boltdb-shipper-compactor
ingester:
  chunk_block_size: 262144
  chunk_idle_period: 3m
  chunk_retain_period: 1m
  lifecycler:
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
  max_transfer_retries: 0
  wal:
    dir: /data/loki/wal
limits_config:
  enforce_metric_name: false
  reject_old_samples: true
  reject_old_samples_max_age: 168h
schema_config:
  configs:
  - from: "2020-10-24"
    index:
      period: 24h
      prefix: index_
    object_store: filesystem
    schema: v11
    store: boltdb-shipper
server:
  http_listen_port: 3100
storage_config:
  boltdb_shipper:
    active_index_directory: /data/loki/boltdb-shipper-active
    cache_location: /data/loki/boltdb-shipper-cache
    cache_ttl: 24h
    shared_store: filesystem
  filesystem:
    directory: /data/loki/chunks
table_manager:
  retention_deletes_enabled: false
  retention_period: 0s
`

const lokiSimpleConfigTplForSimpleScale = `
auth_enabled: false
common:
  path_prefix: /var/loki
  replication_factor: 1
  ring:
    kvstore:
      store: memberlist
  storage:
    filesystem:
      chunks_directory: /var/loki/chunks
      rules_directory: /var/loki/rules
limits_config:
  enforce_metric_name: false
  max_cache_freshness_per_query: 10m
  reject_old_samples: true
  reject_old_samples_max_age: 168h
memberlist:
  join_members:
  - 'loki-ss-loki-simple-scalable-memberlist'
schema_config:
  configs:
  - from: "2020-09-07"
    index:
      period: 24h
      prefix: loki_index_
    object_store: filesystem
    schema: v11
    store: boltdb-shipper
server:
  http_listen_port: 3100
`

// loki configmap
func NewLokiConfigMap(loki *crapiv1alpha1.Loki, mod string) (*apicorev1.ConfigMap, error) {
	var t *template.Template
	switch mod {
	case string(crapiv1alpha1.ModeKindMonolithic):
		t = template.Must(template.New("loki").Parse(lokiSimpleConfigTplForMonolithic))
	case string(crapiv1alpha1.ModeKindSampleScalable):
		t = template.Must(template.New("loki").Parse(lokiSimpleConfigTplForSimpleScale))
	}
	var cfgRaw bytes.Buffer
	if err := t.Execute(&cfgRaw, loki.Spec.Config); err != nil {
		return nil, err
	}
	cm := &apicorev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getLokiConfigMapName(loki, mod),
			Namespace: loki.GetNamespace(),
			OwnerReferences: getResourceOwnerReference(loki),
		},
		Data: map[string]string{"loki.yaml": cfgRaw.String()},
	}
	return cm, nil
}



