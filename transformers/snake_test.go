// Copyright (c) 2020 Brandon Buck

package transformers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/bbuck/luna/transformers"
)

var _ = Describe("StringToSnake", func() {
	DescribeTable("when called",
		func(input, expected string) {
			Ω(transformers.StringToSnake(input)).Should(Equal(expected))
		},
		Entry("basic exported name", "HelloWorld", "hello_world"),
		Entry("camel cased names", "helloWorld", "hello_world"),
		Entry("names with acryonyms", "HelloHTML", "hello_html"),
		Entry("names with acryonyms and more", "HelloHTMLStuff", "hello_html_stuff"),
	)
})
