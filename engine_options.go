// Copyright (c) 2020 Brandon Buck

package luna

// NamingConvention defines how Go names should be converted into Lua names when
// passing values into the Engine.
type NamingConvention int8

const (
	// SnakeCaseAndPascalCase converts all Go names to both their snake_case
	// and Pascal case (ex for 'HelloWorld' you get 'hello_world' and
	// 'HelloWorld')
	SnakeCaseAndPascalCase NamingConvention = iota

	// SnakeCase converts Go names into snake_case only (Default)
	SnakeCase

	// PascalCase converts Go names into Go-exported type case normally
	// (essentially meaning the exported name is unchanged when transitioning
	// to Lua) only.
	PascalCase

	// CamelCase converts Go names into camelCased only.
	CamelCase
)

// EngineOptions allows for customization of a lua.Engine such as altering
// the names of fields and methods as well as whether or not to open all
// libraries.
type EngineOptions struct {
	// OpenLibs determines if the engine should enable the core Lua libraries
	// when creating a new State, if you're looking for security you should not
	// enable this.
	OpenLibs bool

	// FieldCasing defines how the name of a Go struct field should be converted
	// when being passed to Lua.
	FieldCasing NamingConvention

	// MethodCasing defines how the name of a Go struct/interface method should
	// be converted when being passed to Lua.
	MethodCasing NamingConvention
}
