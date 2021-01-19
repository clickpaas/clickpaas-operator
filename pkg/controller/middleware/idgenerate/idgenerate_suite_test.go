package idgenerate_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIdgenerate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Idgenerate Suite")
}
