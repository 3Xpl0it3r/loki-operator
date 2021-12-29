package loki_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLoki(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Loki Suite")
}
