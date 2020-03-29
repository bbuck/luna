// Copyright (c) 2020 Brandon Buck

package luna

import (
	"reflect"

	"github.com/bbuck/luna/transformers"
)

func toSnake(t reflect.Type, s reflect.StructField) []string {
	return []string{transformers.StringToSnake(s.Name)}
}

func toPascal(t reflect.Type, s reflect.StructField) []string {
	return []string{s.Name}
}

func toCamel(t reflect.Type, s reflect.StructField) []string {
	return []string{transformers.StringToCamel(s.Name)}
}
