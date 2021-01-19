package diamond_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDiamond(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Diamond Suite")
}
