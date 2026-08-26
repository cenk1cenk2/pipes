package write

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWrite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Select Environment Write Suite")
}
