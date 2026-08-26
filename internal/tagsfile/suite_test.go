package tagsfile_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTagsFile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tags File Suite")
}
