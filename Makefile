test:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega/...
	ginkgo pkg/controller/
	ginkgo pkg/controller/loki
	ginkgo pkg/controller/promtail
	ginkgo pkg/operator
	ginkgo pkg/operator/loki
	ginkgo pkg/operator/promtail
