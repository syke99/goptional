package goptional

import (
	"encoding/json"
	"reflect"
)

type Goptional[T any] interface {
	Exists(fn func(T)) Goptional[T]
	Val() T
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type goption[T any] struct {
	ptr *T
}

// NewGoptional takes a pointer to the variable
// you want to make Optional and returns it as
// a Goptional[T]
func NewGoptional[T any](opt T) Goptional[T] {
	if reflect.TypeOf(opt).Kind() != reflect.Pointer {
		return nil
	}

	if reflect.ValueOf(opt).IsNil() {
		return nil
	}

	return &goption[T]{
		ptr: &opt,
	}
}

// Exists handles a nil check on your Goptional
// variable and if it isn't nil, passes the underlying
// value to fn
func (g *goption[T]) Exists(fn func(T)) Goptional[T] {
	if !isNil(g.ptr) {
		fn(*g.ptr)
	}
	return g
}

// Val returns the underlying variable if it exists,
// or nil if it doesn't
func (g *goption[T]) Val() T {
	if !isNil(g.ptr) {
		return *g.ptr
	}
	return *(*T)(nil)
}

func isNil[T any](v T) bool {
	return reflect.ValueOf(v).IsNil()
}

// MarshalJSON allows Goptionals to safely implement
// the json.Marshaler interface. This allows fields
// in a struct that will be de/serialized to/from
// JSON to be of type Goptional
func (g *goption[T]) MarshalJSON() ([]byte, error) {
	if !isNil(g.ptr) {
		return json.Marshal(*g.ptr)
	}
	return nil, nil
}

// UnmarshalJSON allows Goptionals to safely implement
// the json.Unmarshaler interface. This allows fields
// in a struct that will be de/serialized to/from
// JSON to be of type Goptional
func (g *goption[T]) UnmarshalJSON(data []byte) error {
	if !isNil(g.ptr) {
		return json.Unmarshal(data, g.ptr)
	}
	return nil
}
