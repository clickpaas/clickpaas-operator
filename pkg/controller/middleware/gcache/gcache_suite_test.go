package gcache_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGcache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gcache Suite")
}
