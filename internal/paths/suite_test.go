package paths_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPaths(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Paths Suite")
}
