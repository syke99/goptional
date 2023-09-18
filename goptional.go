package goptional

import (
	"encoding/json"
	"reflect"
)

type Goptional[T any] interface {
	Exists(fn func(T)) Goptional[T]
	ExistsElse(fn func(T), el func() T) Goptional[T]
	Val() T
	ValOr(or T) T
	ValElse(fn func() T) T
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type goption[T any] struct {
	ptr     *T
	present bool
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
		ptr:     &opt,
		present: true,
	}
}

// Exists handles a nil check on the underlying variable
// of your Goptional and if it isn't nil, passes the underlying
// value to fn
func (g *goption[T]) Exists(fn func(T)) Goptional[T] {
	if g.present {
		fn(*g.ptr)
	}
	return g
}

// ExistsElse handles a nil check on the underlying variable
// of your Goptional and if it isn't nil, passes the underlying
// variable to fn. If it is nil, it calls el, sets the underlying
// variable to the result, then passes that same underlying value
// to fn
func (g *goption[T]) ExistsElse(fn func(T), el func() T) Goptional[T] {
	if g.present {
		fn(*g.ptr)
	} else {
		t := el()
		g.ptr = &t
		fn(*g.ptr)
	}
	return g
}

// Val returns the underlying variable if it exists,
// or nil if it doesn't
func (g *goption[T]) Val() T {
	var t T
	if g.present {
		return *g.ptr
	}
	return t
}

// ValOr returns the underlying variable if it exists,
// or or if it does not
func (g *goption[T]) ValOr(or T) T {
	if g.present {
		return *g.ptr
	}
	return or
}

// ValElse returns the underlying variable if it exists,
// or calls fn and sets the result as the new underlying
// variable and returns it
func (g *goption[T]) ValElse(fn func() T) T {
	if g.present {
		return *g.ptr
	} else {
		t := fn()
		g.ptr = &t
	}
	return *g.ptr
}

// MarshalJSON allows Goptionals to safely implement
// the json.Marshaler interface. This allows fields
// in a struct that will be de/serialized to/from
// JSON to be of type Goptional
func (g *goption[T]) MarshalJSON() ([]byte, error) {
	if g.present {
		return json.Marshal(*g.ptr)
	}
	return nil, nil
}

// UnmarshalJSON allows Goptionals to safely implement
// the json.Unmarshaler interface. This allows fields
// in a struct that will be de/serialized to/from
// JSON to be of type Goptional
func (g *goption[T]) UnmarshalJSON(data []byte) error {
	if g.present {
		return json.Unmarshal(data, g.ptr)
	}
	return nil
}
