// Copyright (c) 2020 Brandon Buck

package luna_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLuna(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Luna Suite")
}

var (
	fileName     = "test_lua.lua"
	fileContents = `
	function give_me_one()
  		return 1
	end
	`
)

var _ = BeforeSuite(func() {
	file, err := os.Create(fileName)
	if err != nil {
		Fail(err.Error())
	}
	defer file.Close()
	fmt.Fprintln(file, fileContents)
})

var _ = AfterSuite(func() {
	err := os.Remove(fileName)
	if err != nil {
		Fail(err.Error())
	}
})
