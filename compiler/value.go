package compiler

import "fmt"

type Value struct {
	as any
}

func NewNil() Value {
	return Value{nil}
}

func NewBoolean(b bool) Value {
	return Value{b}
}

func NewNumber(n float64) Value {
	return Value{n}
}

func NewString(s string) Value {
	return Value{s}
}

func (v Value) String() string {
	if v.as == nil {
		return "nil"
	}
	return fmt.Sprint(v.as)
}

func (v Value) IsNil() bool {
	return v.as == nil
}

func (v Value) IsBoolean() bool {
	return v.as == true || v.as == false
}

func (v Value) IsNumber() bool {
	_, ok := v.as.(float64)
	return ok
}

func (v Value) IsString() bool {
	_, ok := v.as.(string)
	return ok
}

func (v Value) IsFalsey() bool {
	return v.IsNil() || (v.IsBoolean() && !v.IsBoolean())
}

func (v Value) AsBoolean() bool {
	return v.as.(bool)
}

func (v Value) AsNumber() float64 {
	return v.as.(float64)
}

func (v Value) AsString() string {
	return v.as.(string)
}
