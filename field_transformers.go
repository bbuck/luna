// Copyright (c) 2020 Brandon Buck

package luna

import (
	"reflect"

	"github.com/bbuck/luna/transformers"
)

// FieldTransformer is a function that takes information about a struct field and
// returns an array of names to use for this field in Lua.
type FieldTransformer func(reflect.Type, reflect.StructField) []string

func fieldToSnake(t reflect.Type, s reflect.StructField) []string {
	return []string{transformers.StringToSnake(s.Name)}
}

func fieldToPascal(t reflect.Type, s reflect.StructField) []string {
	return []string{s.Name}
}

func fieldToCamel(t reflect.Type, s reflect.StructField) []string {
	return []string{transformers.StringToCamel(s.Name)}
}

// MethodTransformer is a function that takes information about a struct method
// and returns an array of names to use for this method in Lua.
type MethodTransformer func(reflect.Type, reflect.Method) []string

func methodToSnake(t reflect.Type, s reflect.Method) []string {
	return []string{transformers.StringToSnake(s.Name)}
}

func methodToPascal(t reflect.Type, s reflect.Method) []string {
	return []string{s.Name}
}

func methodToCamel(t reflect.Type, s reflect.Method) []string {
	return []string{transformers.StringToCamel(s.Name)}
}
