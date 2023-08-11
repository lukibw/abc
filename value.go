package abc

import "fmt"

type value struct {
	as any
}

func newNil() value {
	return value{nil}
}

func newBoolean(b bool) value {
	return value{b}
}

func newNumber(n float64) value {
	return value{n}
}

func (v value) String() string {
	if v.as == nil {
		return "nil"
	}
	return fmt.Sprint(v.as)
}

func (v value) kind() string {
	switch v.as.(type) {
	case nil:
		return "nil"
	case bool:
		return "boolean"
	case float64:
		return "number"
	default:
		panic(fmt.Sprintf("value has unexpected type '%T'", v.as))
	}
}

func (v value) isNil() bool {
	return v.as == nil
}

func (v value) isBoolean() bool {
	return v.as == true || v.as == false
}

func (v value) isNumber() bool {
	_, ok := v.as.(float64)
	return ok
}

func (v value) asBoolean() bool {
	return v.as.(bool)
}

func (v value) asNumber() float64 {
	return v.as.(float64)
}
