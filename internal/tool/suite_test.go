package tool_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTool(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tool Suite")
}
