package generate

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGenerate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pulumi Generate Suite")
}
