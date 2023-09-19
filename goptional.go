package goptional

import "encoding/json"

type Goptional[T any] interface {
	Exists(fn func(*T) error) error
	ExistsNil() Goptional[T]
	Map(fn func(*T)) Goptional[T]
	FlatMap(fn func(*T)) Goptional[T]
	MapElse(fn func(*T), el func() T) Goptional[T]
	FlatMapElse(fn func(*T), el func() T) Goptional[T]
	Val() any
	ValOr(or T) any
	ValElse(fn func() T) any
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
	unwrapVal() *T
	isWrapped() bool
}

type goption[T any] struct {
	ptr     any
	present bool
	wrapped bool
}

// NewGoptional takes a pointer to the variable you want to
// make optional and returns it as a Goptional[T]
func NewGoptional[T any](opt *T) Goptional[T] {
	if opt == nil {
		return &goption[T]{
			ptr:     nil,
			present: false,
			wrapped: false,
		}
	}

	return &goption[T]{
		ptr:     opt,
		present: true,
		wrapped: false,
	}
}

func Wrap[T any](g Goptional[T]) Goptional[T] {
	if g.Val() == nil {
		return g
	}

	return &goption[T]{
		ptr:     &g,
		present: true,
		wrapped: true,
	}
}

func Unwrap[T any](g Goptional[T]) *T {
	return g.unwrapVal()
}

// Exists checks if the underlying value stored in your Goptional
// is present or not, and if so, passes that value to fn. Otherwise,
func (g *goption[T]) Exists(fn func(*T) error) error {
	if g.ptr != nil {
		return fn(g.ptr.(*T))
	}
	return nil
}

func (g *goption[T]) ExistsNil() Goptional[T] {
	if g.ptr != nil {
		g.ptr = nil
		g.present = false
	}
	return nil
}

// Map handles a nil check on the underlying variable
// of your Goptional and if it isn't nil, passes the underlying
// value to fn
func (g *goption[T]) Map(fn func(*T)) Goptional[T] {
	if g.present && g.ptr != nil {
		fn(g.ptr.(*T))
	}
	return g
}

// FlatMap is like Map, but flattens the Goptional is has
// been called on if it has been wrapped

func (g *goption[T]) FlatMap(fn func(*T)) Goptional[T] {
	var v *T

	if g.wrapped {
		v = g.unwrapVal()
	} else if g.present && g.ptr != nil {
		v = g.ptr.(*T)
	}

	if v != nil {
		fn(v)
	}

	return g
}

// MapElse handles a nil check on the underlying variable
// of your Goptional and if it isn't nil, passes the underlying
// variable to fn. If it is nil, it calls el, sets the underlying
// variable to the result, then passes that same underlying value
// to fn

func (g *goption[T]) MapElse(fn func(*T), el func() T) Goptional[T] {
	if g.present && g.ptr != nil {
		fn(g.ptr.(*T))
	} else {
		v := el()
		fn(&v)
		g.ptr = &v
		g.present = true
	}
	return g
}

// FlatMapElse is like MapElse, but flattens the Goptional is has
// been called on if it has been wrapped
func (g *goption[T]) FlatMapElse(fn func(*T), el func() T) Goptional[T] {
	var v *T

	if g.wrapped {
		v = g.unwrapVal()
	} else if g.present && g.ptr != nil {
		v = g.ptr.(*T)
	} else {
		t := el()
		v = &t
		g.ptr = &t
		g.present = true
	}

	fn(v)

	return g
}

// Val returns a pointer to the underlying variable if it exists,
// or nil if it doesn't
func (g *goption[T]) Val() any {
	if g.present && g.ptr != nil {
		return g.ptr
	}
	return nil
}

func (g *goption[T]) unwrapVal() *T {
	var val *T
	if g.wrapped {
		x := *g.ptr.(*Goptional[T])
		return x.unwrapVal()
	}

	if g.present && g.ptr != nil {
		return g.ptr.(*T)
	}
	return val
}

func (g *goption[T]) isWrapped() bool {
	return g.wrapped
}

// ValOr returns the underlying variable if it exists,
// or or if it does not
func (g *goption[T]) ValOr(or T) any {
	if g.present && g.ptr != nil {
		return g.ptr.(*T)
	}
	return or
}

// ValElse returns the underlying variable if it exists,
// or calls fn and sets the result as the new underlying
// variable and returns it
func (g *goption[T]) ValElse(fn func() T) any {
	if g.present && g.ptr != nil {
		return g.ptr.(*T)
	}
	t := fn()
	g.ptr = &t
	g.present = true
	return t
}

// MarshalJSON allows Goptionals to safely implement
// the json.Marshaler interface. This allows fields
// in a struct that will be de/serialized to/from
// JSON to be of type Goptional
func (g *goption[T]) MarshalJSON() ([]byte, error) {
	if g.wrapped {
		v := g.unwrapVal()
		return json.Marshal(v)
	}

	if g.present && g.ptr != nil {
		return json.Marshal(g.ptr)
	}
	return nil, nil
}

// UnmarshalJSON allows Goptionals to safely implement
// the json.Unmarshaler interface. This allows fields
// in a struct that will be de/serialized to/from
// JSON to be of type Goptional
func (g *goption[T]) UnmarshalJSON(data []byte) error {
	if g.wrapped {
		v := g.unwrapVal()
		return json.Unmarshal(data, v)
	}

	if g.present && g.ptr != nil {
		return json.Unmarshal(data, g.ptr)
	}
	return nil
}
