package ci_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CI Suite")
}
