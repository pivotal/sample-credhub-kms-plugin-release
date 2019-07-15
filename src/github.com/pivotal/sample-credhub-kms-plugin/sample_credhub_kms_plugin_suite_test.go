package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSampleCredhubKmsPlugin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SampleCredhubKmsPlugin Suite")
}
