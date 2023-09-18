package goptional

import "reflect"

type Goptional[T any] interface {
	Exists(fn func(T)) Goptional[T]
}

type goption[T any] struct {
	ptr *T
}

func NewGoptional[T any](opt T) Goptional[T] {
	if reflect.TypeOf(opt).Kind() != reflect.Pointer {
		return nil
	}

	return &goption[T]{
		ptr: &opt,
	}
}

func (g *goption[T]) Exists(fn func(T)) Goptional[T] {
	if !isNil(g.ptr) {
		fn(*g.ptr)
	}
	return g
}

func isNil[T any](v T) bool {
	return reflect.ValueOf(v).IsNil()
}
