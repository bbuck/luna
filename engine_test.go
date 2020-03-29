// Copyright (c) 2020 Brandon Buck

package luna_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bbuck/luna"
)

var _ = Describe("LuaEngine", func() {
	var (
		err          error
		engine       *Engine
		stringScript = `
			function hello(name)
				return "Hello, " .. name .. "!"
			end
		`
	)

	BeforeEach(func() {
		engine = NewEngine()
	})

	AfterEach(func() {
		engine.Close()
	})

	Context("when closed", func() {
		BeforeEach(func() {
			engine.Close()
		})

		It("no longer functions", func() {
			Skip("Skipping until fixed")
			_, err := engine.Call("hello", 1, "World")
			Ω(err).ShouldNot(BeNil())
		})
	})

	Context("when loading from a string", func() {
		BeforeEach(func() {
			err = engine.DoString(stringScript)
		})

		It("should not fail", func() {
			Ω(err).Should(BeNil())
		})

		Context("when calling a method", func() {
			var (
				results []*Value
				err     error
			)

			BeforeEach(func() {
				results, err = engine.Call("hello", 1, "World")
			})

			It("does not return an error", func() {
				Ω(err).Should(BeNil())
			})

			It("returns 1 result", func() {
				Ω(len(results)).Should(Equal(1))
			})

			It("doesn't return nil", func() {
				Ω(results[0]).ShouldNot(Equal(engine.Nil()))
			})

			It("returns the string 'Hello, World!'", func() {
				Ω(results[0].AsString()).Should(Equal("Hello, World!"))
			})
		})
	})

	Context("when loading from a file", func() {
		BeforeEach(func() {
			err = engine.DoFile(fileName)
		})

		It("shoult not fail", func() {
			Ω(err).Should(BeNil())
		})

		Context("when calling a method", func() {
			var (
				results []*Value
				err     error
			)

			BeforeEach(func() {
				results, err = engine.Call("give_me_one", 1)
			})

			It("does not return an error", func() {
				Ω(err).Should(BeNil())
			})

			It("return 1 result", func() {
				Ω(len(results)).Should(Equal(1))
			})

			It("does not return nil", func() {
				Ω(results[0]).ShouldNot(Equal(engine.Nil()))
			})

			It("returns the number 1", func() {
				Ω(results[0].AsNumber()).Should(Equal(float64(1)))
			})
		})
	})

	Describe("Call()", func() {
		var (
			results []*Value
			err     error
			script  = `
				function swap(a, b)
					return b, a
				end
			`
			a                float64 = 10.0
			b                float64 = 20.0
			aResult, bResult float64
		)

		BeforeEach(func() {
			engine.DoString(script)
			results, err = engine.Call("swap", 2, a, b)
			if err == nil {
				aResult = results[0].AsNumber()
				bResult = results[1].AsNumber()
			}
		})

		It("does not return an error", func() {
			Ω(err).Should(BeNil())
		})

		It("returns two results", func() {
			Ω(len(results)).Should(Equal(2))
		})

		It("returns the second input first", func() {
			Ω(aResult).Should(Equal(b))
		})

		It("returns the first input second", func() {
			Ω(bResult).Should(Equal(a))
		})
	})

	Describe("SetGlobal()", func() {
		var (
			results []*Value
			err     error
		)

		BeforeEach(func() {
			engine.SetGlobal("gbl", "testing")
			err = engine.DoString(`
			function get_gbl()
				return gbl
			end
			`)
			if err != nil {
				Fail(err.Error())
			}
			results, err = engine.Call("get_gbl", 1)
		})

		It("does not fail", func() {
			Ω(err).Should(BeNil())
		})

		It("returns one value", func() {
			Ω(len(results)).Should(Equal(1))
		})

		It("returns the value assigned to the global", func() {
			Ω(results[0].AsString()).Should(Equal("testing"))
		})
	})

	Describe("GetGlobal()", func() {
		var (
			value *Value
			err   error
		)

		BeforeEach(func() {
			err = engine.DoString(`
				word = "testing"
			`)
			if err != nil {
				Fail(err.Error())
			}
			value = engine.GetGlobal("word")
		})

		It("doesn't return nil", func() {
			Ω(value).ShouldNot(Equal(engine.Nil()))
		})

		It("returns the correct string", func() {
			Ω(value.AsString()).Should(Equal("testing"))
		})
	})

	Describe("RegisterFunc()", func() {
		Context("when registering a raw Go function", func() {
			var (
				results []*Value
				err     error
				called  bool
			)

			BeforeEach(func() {
				engine.RegisterFunc("add", func(a, b int) int {
					called = true
					return a + b
				})
				results, err = engine.Call("add", 1, 10, 11)
			})

			It("should no fail", func() {
				Ω(err).Should(BeNil())
			})

			It("marks the called variable", func() {
				Ω(called).Should(BeTrue())
			})

			It("does not return nil", func() {
				Ω(results[0]).ShouldNot(Equal(engine.Nil()))
			})

			It("returns 1 value", func() {
				Ω(len(results)).Should(Equal(1))
			})

			It("returns a value that passed through the Go function", func() {
				Ω(results[0].AsNumber()).Should(Equal(float64(21)))
			})
		})

		Context("when registering a lua specific function", func() {
			var (
				results []*Value
				err     error
				called  bool
			)

			BeforeEach(func() {
				engine.RegisterFunc("sub", func(e *Engine) int {
					second := e.PopInt64()
					first := e.PopInt64()

					if first == 11 && second == 10 {
						called = true
					}

					e.PushValue(first - second)

					return 1
				})
				results, err = engine.Call("sub", 1, 11, 10)
			})

			It("does not fail", func() {
				Ω(err).Should(BeNil())
			})

			It("returns 1 value", func() {
				Ω(len(results)).Should(Equal(1))
			})

			It("marks the variable called", func() {
				Ω(called).Should(BeTrue())
			})

			It("does not return nil", func() {
				Ω(results[0]).ShouldNot(Equal(engine.Nil()))
			})

			It("returns the correct value", func() {
				Ω(results[0].AsNumber()).Should(Equal(float64(1)))
			})
		})
	})

	Describe("passing in go objects", func() {
		var obj = TestObject{}

		BeforeEach(func() {
			engine.DoString(`
				function call_by_value_fn(obj)
				  return obj:GetStringFromValue()
				end

				function call_by_ptr_fn(obj)
					return obj:GetStringFromPtr()
				end
			`)
		})

		Context("calling methods by value", func() {
			var (
				result []*Value
				cerr   error
			)

			BeforeEach(func() {
				result, cerr = engine.Call("call_by_value_fn", 1, obj)
			})

			It("should not fail", func() {
				Ω(cerr).Should(BeNil())
			})

			It("should return the correct value", func() {
				Ω(len(result)).Should(BeNumerically(">", 0))
				Ω(result[0].AsString()).Should(Equal("success"))
			})
		})

		Context("calling methods by pointer", func() {
			var (
				result []*Value
				cerr   error
			)

			BeforeEach(func() {
				result, cerr = engine.Call("call_by_ptr_fn", 1, &obj)
			})

			It("should not fail", func() {
				Ω(cerr).Should(BeNil())
			})

			It("should return the correct value", func() {
				Ω(len(result)).Should(BeNumerically(">", 0))
				Ω(result[0].AsString()).Should(Equal("success"))
			})
		})
	})

	Describe("using table generators", func() {
		var (
			table          *Value
			results        []*Value
			errOne, errTwo error
			one            *Value
			two            *Value
		)

		BeforeEach(func() {
			engine.DoString(`
                function getValueAtKey(tbl, key)
                    return tbl[key]
                end
            `)
		})

		Context("ValueFromMap", func() {
			m := map[string]interface{}{
				"one": 2,
				"two": "too",
			}

			BeforeEach(func() {
				table = engine.TableFromMap(m)
				results, errOne = engine.Call("getValueAtKey", 1, table, "one")
				if len(results) > 0 {
					one = results[0]
				}
				results, errTwo = engine.Call("getValueAtKey", 1, table, "two")
				if len(results) > 0 {
					two = results[0]
				}
			})

			It("didn't fail to fetch 'one'", func() {
				Ω(errOne).Should(BeNil())
			})

			It("fetched a number", func() {
				Ω(one.IsNumber()).Should(BeTrue())
			})

			It("fetch the number 2", func() {
				Ω(one.AsNumber()).Should(Equal(float64(2)))
			})

			It("didn't fail to fetch 'two'", func() {
				Ω(errTwo).Should(BeNil())
			})

			It("fetched a string", func() {
				Ω(two.IsString()).Should(BeTrue())
			})

			It("fetch the string 'too'", func() {
				Ω(two.AsString()).Should(Equal("too"))
			})
		})

		Context("ValueFromSlice", func() {
			s := []int{1, 2, 3}

			BeforeEach(func() {
				table = engine.TableFromSlice(s)
				results, errOne = engine.Call("getValueAtKey", 1, table, 1)
				if len(results) > 0 {
					one = results[0]
				}
				results, errTwo = engine.Call("getValueAtKey", 1, table, 2)
				if len(results) > 0 {
					two = results[0]
				}
			})

			It("has 3 values", func() {
				Ω(table.Len()).Should(Equal(3))
			})

			It("didn't fail to fetch #1", func() {
				Ω(errOne).Should(BeNil())
			})

			It("fetched a number", func() {
				Ω(one.IsNumber()).Should(BeTrue())
			})

			It("fetch the number 1", func() {
				Ω(one.AsNumber()).Should(Equal(float64(1)))
			})

			It("didn't fail to fetch #2", func() {
				Ω(errTwo).Should(BeNil())
			})

			It("fetched a number", func() {
				Ω(two.IsNumber()).Should(BeTrue())
			})

			It("fetch the number 2", func() {
				Ω(two.AsNumber()).Should(Equal(float64(2)))
			})
		})
	})
})

type TestObject struct{}

func (t TestObject) GetStringFromValue() string {
	return "success"
}

func (t *TestObject) GetStringFromPtr() string {
	return "success"
}
