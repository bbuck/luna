// Copyright (c) 2020 Brandon Buck

package transformers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/bbuck/luna/transformers"
)

var _ = Describe("StringToCamel", func() {
	DescribeTable("when called",
		func(input, expected string) {
			Ω(transformers.StringToCamel(input)).Should(Equal(expected))
		},
		Entry("basic exported string", "HelloWorld", "helloWorld"),
		Entry("camel cased string", "helloWorld", "helloWorld"),
		Entry("string with acronyms", "HelloHTML", "helloHTML"),
		Entry("string with acronyms and more", "HelloHTMLStuff", "helloHTMLStuff"),
		Entry("string that starts with acronym", "HTMLLibrary", "htmlLibrary"),
	)
})
