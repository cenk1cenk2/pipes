package gitlab

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGitlab(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gitlab Suite")
}
