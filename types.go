// Copyright (c) 2020 Brandon Buck

package luna

// ScriptFunction is a Go function intended to be called from within Lua code.
// It receives the Engine that it was called from within as its argument which
// can be used to extract arguments from the Lua stack. It is expected to push
// all return values back onto the Lua stack and the function returns the number
// of return values it pushed.
type ScriptFunction func(*Engine) int

// TableMap is simple type definition around map[string]interface{} which is a
// more clear and friendly Go type to use when defining a Table meant to be sent
// into Lua.
type TableMap map[string]interface{}
