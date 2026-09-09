package publish

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPublish(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Helm Publish Suite")
}
