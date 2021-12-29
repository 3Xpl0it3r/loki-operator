package promtail_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPromtail(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Promtail Suite")
}
