package cobertura_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCobertura(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cobertura Report Suite")
}
