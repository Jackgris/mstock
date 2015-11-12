package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMstock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mstock Suite")
}
