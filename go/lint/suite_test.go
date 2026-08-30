package lint

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLint(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Go Lint Suite")
}
