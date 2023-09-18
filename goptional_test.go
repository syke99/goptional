package goptional

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func transform(val *testType) {
	val.greeting = "hello"
}

type testType struct {
	greeting string
}

func TestNewGoptionalPointer(t *testing.T) {
	// Arrange
	tt := testType{}

	// Act
	opt := NewGoptional(&tt)

	// Assert
	assert.NotNil(t, opt)
}

func TestExists(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	// Act
	opt.Exists(transform)

	// Assert
	assert.Equal(t, "hello", tt.greeting)
}

func TestExistsElseExists(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	// Act
	opt.ExistsElse(transform, func() *testType {
		return &testType{}
	})

	// Assert
	assert.Equal(t, "hello", tt.greeting)
}

func TestExistsElseDoesntExists(t *testing.T) {
	// Arrange
	opt := NewGoptional[*testType](nil)

	// Act
	opt.ExistsElse(transform, func() *testType {
		return &testType{}
	})

	// Assert
	assert.Equal(t, "hello", any(opt.Val()).(*testType).greeting)
}

func TestVal(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	// Act
	v := opt.Val()

	// Assert
	assert.NotNil(t, v)
}

func TestValNoVal(t *testing.T) {
	// Arrange
	opt := NewGoptional[*testType](nil)

	// Act
	v := opt.Val()

	// Assert
	assert.Nil(t, v)
}

func TestValOr(t *testing.T) {
	// Arrange
	tt := testType{greeting: "hello"}

	opt := NewGoptional[*testType](nil)

	// Act
	v := opt.ValOr(&tt)

	// Assert
	assert.Equal(t, "hello", any(v).(*testType).greeting)
}

func TestValElse(t *testing.T) {
	// Arrange
	opt := NewGoptional[*testType](nil)

	// Act
	v := opt.ValElse(func() *testType {
		return &testType{greeting: "hello"}
	})

	// Assert
	assert.Equal(t, "hello", any(v).(*testType).greeting)
}
