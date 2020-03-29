// Copyright (c) 2020 Brandon Buck

package luna

import (
	"bytes"
	"fmt"
	"math"

	lua "github.com/yuin/gopher-lua"
)

// Inspecter defines an type that can respond to the Inspect function. This is
// similar to fmt.Stringer in that it's a method that returns a string, but the
// goal with Inspecter over fmt.Stringer is to provide debug information in
// the string output rather than (potentially) user facing output.
type Inspecter interface {
	Inspect(string) string
}

// Value is a utility wrapper for lua.LValue that provies conveinient methods
// for casting.
type Value struct {
	lval  lua.LValue
	owner *Engine
}

// String makes Value conform to Stringer
func (v *Value) String() string {
	return v.lval.String()
}

// AsRaw returns the best associated Go type, ingoring functions and any other
// odd types. Only concerns itself with string, bool, nil, number and user data
// types. Tables are again, ignored.
func (v *Value) AsRaw() interface{} {
	switch v.lval.Type() {
	case lua.LTString:
		return v.AsString()
	case lua.LTBool:
		return v.AsBool()
	case lua.LTNil:
		return nil
	case lua.LTNumber:
		return v.AsNumber()
	case lua.LTUserData:
		return v.Interface()
	case lua.LTTable:
		if v.Len() > 0 {
			return v.AsSliceInterface()
		}

		return v.AsMapStringInterface()
	}

	return nil
}

// Inspect is similar to AsString except that it's designed to display values
// for debug purposes.
func (v *Value) Inspect(indent string) string {
	nextIndent := indent + "  "

	switch v.lval.Type() {
	case lua.LTString:
		return fmt.Sprintf("%q", v.AsString())
	case lua.LTBool:
		if v.IsTrue() {
			return "true"
		}

		return "false"
	case lua.LTNil:
		return "nil"
	case lua.LTNumber:
		n := v.AsNumber()
		if math.Floor(n) == n && n >= float64(math.MinInt64) && n <= float64(math.MaxInt64) {
			return fmt.Sprintf("%d", int64(n))
		}

		return fmt.Sprintf("%g", n)
	case lua.LTUserData:
		iface := v.Interface()
		switch it := iface.(type) {
		case Inspecter:
			return it.Inspect(indent)
		case fmt.Stringer:
			return it.String()
		default:
			ud := v.asUserData()
			val := v.owner.ValueFor(ud.Metatable)

			mt := v.owner.MetatableFor(ud.Value)
			vals, err := mt.RawGet("ptr_methods").Invoke("inspect", 1, v, indent)
			if err == nil && len(vals) > 0 {
				return vals[0].AsString()
			}

			vals, err = mt.Invoke("inspect", 1, v, indent)
			if err == nil && len(vals) > 0 {
				return vals[0].AsString()
			}

			vals, err = val.Invoke("inspect", 1, v, indent)
			if err != nil || len(vals) == 0 {
				return fmt.Sprintf("%T(%+v)", iface, iface)
			}

			return vals[0].AsString()
		}
	case lua.LTTable:
		vals, err := v.Invoke("inspect", 1, v)
		if err != nil || len(vals) == 0 {
			buf := new(bytes.Buffer)
			buf.WriteString("{\n")
			v.ForEach(func(key, val *Value) {
				buf.WriteString(nextIndent)
				buf.WriteString(fmt.Sprintf("[%s] = %s", key.Inspect(nextIndent), val.Inspect(nextIndent)))
				buf.WriteString(",\n")
			})
			buf.WriteString(indent)
			buf.WriteString("}")

			return buf.String()
		}

		return vals[0].Inspect(nextIndent)
	case lua.LTFunction:
		return "<function>"
	}

	return "nil"
}

// AsString returns the LValue as a Go string
func (v *Value) AsString() string {
	return lua.LVAsString(v.lval)
}

// AsFloat returns the LValue as a Go float64.
// This method will try to convert the Lua value to a number if possible, if
// not then LuaNumber(0) is returned.
func (v *Value) AsFloat() float64 {
	return float64(lua.LVAsNumber(v.lval))
}

// AsNumber is an alias for AsFloat (Lua calls them "numbers")
func (v *Value) AsNumber() float64 {
	return v.AsFloat()
}

// AsBool returns the Lua boolean representation for an object (this works for
// non bool Values)
func (v *Value) AsBool() bool {
	return lua.LVAsBool(v.lval)
}

// AsMapStringInterface will work on a Lua Table to convert it into a go
// map[string]interface. This method is not safe for cyclic objects! You have
// been warned.
func (v *Value) AsMapStringInterface() map[string]interface{} {
	if v.IsTable() {
		result := make(map[string]interface{})
		v.ForEach(func(key, value *Value) {
			var val interface{} = value.AsRaw()
			if value.IsTable() {
				if value.IsMaybeList() {
					val = value.AsSliceInterface()
				} else {
					val = value.AsMapStringInterface()
				}
			}
			result[key.AsString()] = val
		})

		return result
	}

	return nil
}

// AsSliceInterface will convert the Lua table value to a []interface{},
// extracting Go values were possible and preserving references to tables.
func (v *Value) AsSliceInterface() []interface{} {
	if v.IsTable() {
		s := make([]interface{}, 0)
		len := v.Len()
		for i := 1; i <= len; i++ {
			lv := v.Get(i)
			var val interface{} = lv.AsRaw()
			if lv.IsTable() {
				if lv.IsMaybeList() {
					val = lv.AsSliceInterface()
				} else {
					val = lv.AsMapStringInterface()
				}
			}
			s = append(s, val)
		}

		return s
	}

	return nil
}

// Equals will determine if the *Value is equal to the other value. This also
// verifies they are from the same *lua.Engine as well.
func (v *Value) Equals(o interface{}) bool {
	oval := v.owner.ValueFor(o)

	return oval.owner == v.owner && v.owner.state.Equal(v.lval, oval.lval)
}

// IsNil will only return true if the Value wraps LNil.
func (v *Value) IsNil() bool {
	return v == nil || v.lval == nil || v.lval.Type() == lua.LTNil
}

// IsFalse is similar to AsBool except it returns if the Lua value would be
// considered false in Lua.
func (v *Value) IsFalse() bool {
	return lua.LVIsFalse(v.lval)
}

// IsTrue returns whether or not this is a truthy value or not.
func (v *Value) IsTrue() bool {
	return !v.IsFalse()
}

// The following methods allow for type detection

// IsNumber returns true if the stored value is a numeric value.
func (v *Value) IsNumber() bool {
	return v.lval.Type() == lua.LTNumber
}

// IsBool returns true if the stored value is a boolean value.
func (v *Value) IsBool() bool {
	return v.lval.Type() == lua.LTBool
}

// IsFunction returns true if the stored value is a function.
func (v *Value) IsFunction() bool {
	return v.lval.Type() == lua.LTFunction
}

// IsString returns true if the stored value is a string.
func (v *Value) IsString() bool {
	return v.lval.Type() == lua.LTString
}

// IsTable returns true if the stored value is a table.
func (v *Value) IsTable() bool {
	return v.lval.Type() == lua.LTTable
}

// IsMaybeList will try and determine if the table _might_ be used as a list.
// This basically checks the first index (looking for something at 1) so  it's
// not 100% accurate, hence 'Maybe.' Also a list that starts with 'nil' will
// report as not a list.
func (v *Value) IsMaybeList() bool {
	if v.IsTable() {
		if v.RawGet(1).IsNil() {
			return false
		}

		return true
	}

	return false
}

// The following methods allow LTable values to be modified through Go.

// asTable converts the Value into an LTable.
func (v *Value) asTable() (t *lua.LTable) {
	t, _ = v.lval.(*lua.LTable)

	return
}

// IsUserData returns a bool if the Value is an LUserData
func (v *Value) IsUserData() bool {
	return v.lval.Type() == lua.LTUserData
}

// asUserData converts the Value into an LUserData
func (v *Value) asUserData() (t *lua.LUserData) {
	t, _ = v.lval.(*lua.LUserData)

	return
}

// Append maps to lua.LTable.Append
func (v *Value) Append(value interface{}) {
	if v.IsTable() {
		val := getLValue(v.owner, value)

		t := v.asTable()
		t.Append(val)
	}
}

// ForEach maps to lua.LTable.ForEach
func (v *Value) ForEach(cb func(*Value, *Value)) {
	if v.IsTable() {
		actualCb := func(key lua.LValue, val lua.LValue) {
			cb(v.owner.newValue(key), v.owner.newValue(val))
		}
		t := v.asTable()
		t.ForEach(actualCb)
	}
}

// Insert maps to lua.LTable.Insert
func (v *Value) Insert(i int, value interface{}) {
	if v.IsTable() {
		val := getLValue(v.owner, value)

		t := v.asTable()
		t.Insert(i, val)
	}
}

// Len maps to lua.LTable.Len
func (v *Value) Len() int {
	if v.IsTable() {
		t := v.asTable()

		return t.Len()
	}

	return -1
}

// MaxN maps to lua.LTable.MaxN
func (v *Value) MaxN() int {
	if v.IsTable() {
		t := v.asTable()

		return t.MaxN()
	}

	return 0
}

// Next maps to lua.LTable.Next
func (v *Value) Next(key interface{}) (*Value, *Value) {
	if v.IsTable() {
		val := getLValue(v.owner, key)

		t := v.asTable()
		v1, v2 := t.Next(val)

		return v.owner.newValue(v1), v.owner.newValue(v2)
	}

	return v.owner.Nil(), v.owner.Nil()
}

// Remove maps to lua.LTable.Remove
func (v *Value) Remove(pos int) *Value {
	if v.IsTable() {
		t := v.asTable()
		ret := t.Remove(pos)

		return v.owner.newValue(ret)
	}

	return v.owner.Nil()
}

// Helper method for Set and RawSet
func getLValue(e *Engine, item interface{}) lua.LValue {
	switch val := item.(type) {
	case (*Value):
		return val.lval
	case lua.LValue:
		return val
	}

	if e != nil {
		return e.ValueFor(item).lval
	}

	return lua.LNil
}

// Get returns the value associated with the key given if the LuaValue wraps
// a table.
func (v *Value) Get(key interface{}) *Value {
	if v.IsTable() {
		k := getLValue(v.owner, key)
		val := v.owner.state.GetTable(v.lval, k)

		return v.owner.ValueFor(val)
	}

	return nil
}

// Set will assign the field on the object to the given value.
func (v *Value) Set(goKey interface{}, val interface{}) {
	if v.IsTable() {
		key := v.owner.ValueFor(goKey)

		v.owner.SetField(v, key.AsString(), val)
	}
}

// RawSet bypasses any checks for key existence and sets the value onto the
// table with the given key.
func (v *Value) RawSet(goKey interface{}, val interface{}) {
	if v.IsTable() {
		key := getLValue(v.owner, goKey)
		lval := getLValue(v.owner, val)

		v.asTable().RawSet(key, lval)
	}
}

// RawSetInt sets some value at the given integer index value.
func (v *Value) RawSetInt(i int, val interface{}) {
	if v.IsTable() {
		lval := getLValue(v.owner, val)

		v.asTable().RawSetInt(i, lval)
	}
}

// RawGet fetches data from a table, bypassing __index metamethod.
func (v *Value) RawGet(goKey interface{}) *Value {
	if v.IsTable() {
		key := getLValue(v.owner, goKey)
		ret := v.asTable().RawGet(key)

		return v.owner.ValueFor(ret)
	}

	return v.owner.Nil()
}

// The following provde methods for LUserData

// Interface returns the value of the LUserData
func (v *Value) Interface() interface{} {
	if v.IsUserData() {
		t := v.asUserData()

		return t.Value
	}

	return nil
}

// The following provide LFunction methods on Value

// FuncLocalName is a function that returns the local name of a LFunction type
// if this Value objects holds an LFunction.
func (v *Value) FuncLocalName(regno, pc int) (string, bool) {
	if f, ok := v.lval.(*lua.LFunction); ok {
		return f.LocalName(regno, pc)
	}

	return "", false
}

// Invoke will fetch a funtion value on the table (if we're working with a
// table, and then attempt to invoke it if it's a function.
func (v *Value) Invoke(key interface{}, retCount int, argList ...interface{}) ([]*Value, error) {
	var val *Value
	if v.IsUserData() {
		ud := v.lval.(*lua.LUserData)
		mtbl := v.owner.ValueFor(ud.Metatable)
		val = mtbl.Get(key)
		fmt.Printf("\n\n%+v\n\n", val)
	} else {
		val = v.Get(key)
	}

	if val == nil || val.IsNil() || !val.IsFunction() {
		return nil, fmt.Errorf("value doesn't exist or is not a function")
	}

	return val.Call(retCount, argList...)
}

// Call invokes the LuaValue as a function (if it is one) with similar behavior
// to engine.Call. If you're looking to invoke a function on table, then see
// Value.Invoke
func (v *Value) Call(retCount int, argList ...interface{}) ([]*Value, error) {
	if v.IsFunction() && v.owner != nil {
		p := lua.P{
			Fn:      v.lval,
			NRet:    retCount,
			Protect: true,
		}
		args := make([]lua.LValue, len(argList))
		for i, iface := range argList {
			args[i] = getLValue(v.owner, iface)
		}

		err := v.owner.state.CallByParam(p, args...)
		if err != nil {
			return nil, err
		}

		retVals := make([]*Value, retCount)
		for i := 0; i < retCount; i++ {
			retVals[i] = v.owner.ValueFor(v.owner.state.Get(-1))
		}

		return retVals, nil
	}

	return make([]*Value, 0), nil
}
