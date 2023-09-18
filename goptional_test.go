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

func TestNewGoptionalNotPointer(t *testing.T) {
	// Act
	opt := NewGoptional("hello")

	// Assert
	assert.Nil(t, opt)
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
