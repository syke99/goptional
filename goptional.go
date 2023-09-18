package goptional

import (
	"encoding/json"
)

type Goptional[T any] interface {
	Exists(fn func(T)) Goptional[T]
	ExistsElse(fn func(T), el func() T) Goptional[T]
	Val() T
	ValOr(or T) T
	ValElse(fn func() T) T
	MarshalJSON() ([]byte, error)
	//UnmarshalJSON([]byte) error
	unwrapVal() T
	isWrapped() bool
}

type goption[T any] struct {
	ptr     *any
	present bool
	wrapped bool
}

// NewGoptional takes a pointer to the variable you want to
// make optional and returns it as a Goptional[T]
func NewGoptional[T comparable](opt T) Goptional[T] {
	gop := &goption[T]{}

	v := any(opt)

	var isNil T
	if opt == isNil {
		gop.ptr = &v
		gop.present = false
		return gop
	}

	gop.ptr = &v
	gop.present = true

	return gop
}

func Wrap[T comparable](g Goptional[T]) Goptional[T] {
	var n T
	if g.Val() == n {
		return g
	}

	v := any(g)

	return &goption[T]{
		ptr:     &v,
		present: true,
		wrapped: true,
	}
}

func Unwrap[T comparable](g Goptional[T]) T {
	return g.unwrapVal()
}

// Exists handles a nil check on the underlying variable
// of your Goptional and if it isn't nil, passes the underlying
// value to fn
func (g *goption[T]) Exists(fn func(T)) Goptional[T] {
	if g.present {
		v := *g.ptr
		fn(v.(T))
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
		v := *g.ptr
		fn(v.(T))
	} else {
		t := any(el())
		g.ptr = &t
		g.present = true
		v := *g.ptr
		fn(v.(T))
	}
	return g
}

// Val returns the underlying variable if it exists,
// or nil if it doesn't
func (g *goption[T]) Val() T {
	var t T
	if g.present {
		v := *g.ptr
		if g.isWrapped() {
			x := v.(*goption[T])
			unwrapped := *x.ptr
			return unwrapped.(T)
		}
		return v.(T)
	}
	return t
}

func (g *goption[T]) isWrapped() bool {
	return g.wrapped
}

func (g *goption[T]) unwrapVal() T {
	var val T
	if g.isWrapped() {
		v := *g.ptr
		x := v.(*goption[T])
		return x.unwrapVal()
	}

	if g.present {
		v := *g.ptr
		return v.(T)
	}
	return val
}

// ValOr returns the underlying variable if it exists,
// or or if it does not
func (g *goption[T]) ValOr(or T) T {
	if g.present {
		v := *g.ptr
		return v.(T)
	}
	return or
}

// ValElse returns the underlying variable if it exists,
// or calls fn and sets the result as the new underlying
// variable and returns it
func (g *goption[T]) ValElse(fn func() T) T {
	if g.present {
		v := *g.ptr
		return v.(T)
	} else {
		t := any(fn())
		g.ptr = &t
		g.present = true
	}
	v := *g.ptr
	return v.(T)
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

//
//// UnmarshalJSON allows Goptionals to safely implement
//// the json.Unmarshaler interface. This allows fields
//// in a struct that will be de/serialized to/from
//// JSON to be of type Goptional
//func (g *goption[T]) UnmarshalJSON(data []byte) error {
//	if g.isWrapped() {
//		v := g.unwrapVal()
//		return v.UnmarshalJSON(data)
//	}
//
//	if g.present {
//		return json.Unmarshal(data, g.ptr)
//	}
//	return nil
//}
