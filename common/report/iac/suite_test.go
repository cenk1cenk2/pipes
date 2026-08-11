package iac_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIac(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IaC Report Suite")
}
