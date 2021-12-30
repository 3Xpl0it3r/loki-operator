package v1alpha1

func WithDefaultsLoki(loki *Loki) {

}

func WithDefaultsPromtail(promtail *Promtail) {
	// config相关的配置
	if promtail.Spec.ConfigMap != "" {
		return
	}
	if promtail.Spec.ConfigMap == "" {
		if promtail.Spec.Config.Clients.URL == "" {
			promtail.Spec.Config.Clients.URL = defaultLokiGwForPromtail(promtail.GetName())
		}
		if promtail.Spec.Config.Server.HttpListenAddress == "" {
			promtail.Spec.Config.Server.HttpListenAddress = ":9090"
		}
	}
}

func defaultLokiGwForPromtail(promtail string) string {
	return promtail + "-loki-gateway"
}
