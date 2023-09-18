package goptional

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func transform(val *testType) {
	val.Greeting = "hello"
}

type testType struct {
	Greeting string `json:"greeting"`
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
	assert.Equal(t, "hello", tt.Greeting)
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
	assert.Equal(t, "hello", tt.Greeting)
}

func TestExistsElseDoesntExists(t *testing.T) {
	// Arrange
	opt := NewGoptional[*testType](nil)

	// Act
	opt.ExistsElse(transform, func() *testType {
		return &testType{}
	})

	// Assert
	assert.Equal(t, "hello", any(opt.Val()).(*testType).Greeting)
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
	tt := testType{Greeting: "hello"}

	opt := NewGoptional[*testType](nil)

	// Act
	v := opt.ValOr(&tt)

	// Assert
	assert.Equal(t, any(v).(*testType), &tt)
	assert.Equal(t, "hello", any(v).(*testType).Greeting)
}

func TestValElse(t *testing.T) {
	// Arrange
	opt := NewGoptional[*testType](nil)

	// Act
	v := opt.ValElse(func() *testType {
		return &testType{Greeting: "hello"}
	})

	// Assert
	assert.Equal(t, "hello", any(v).(*testType).Greeting)
}

func TestMarshalJSON(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	opt.Exists(transform)

	// Act
	jsonBytes, err := opt.MarshalJSON()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "{\"greeting\":\"hello\"}", string(jsonBytes))
}

func TestMarshalJSONDoesntExist(t *testing.T) {
	// Arrange
	opt := NewGoptional[*testType](nil)

	// Act
	jsonBytes, err := opt.MarshalJSON()

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, jsonBytes)
}

//
//func TestUnmarshalJSON(t *testing.T) {
//	// Arrange
//	tt := testType{}
//
//	opt := NewGoptional(&tt)
//
//	jsonString := "{\"greeting\":\"hello\"}"
//
//	// Act
//	err := opt.UnmarshalJSON([]byte(jsonString))
//
//	// Assert
//	assert.NoError(t, err)
//	assert.Equal(t, "hello", any(opt.Val()).(*testType).Greeting)
//}
//
//func TestUnmarshalJSONDoesntExist(t *testing.T) {
//	// Arrange
//	opt := NewGoptional[*testType](nil)
//
//	jsonString := "{\"greeting\":\"hello\"}"
//
//	// Act
//	err := opt.UnmarshalJSON([]byte(jsonString))
//
//	// Assert
//	assert.NoError(t, err)
//}

func TestWrap(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	// Act
	opt2 := Wrap(opt)

	// Assert
	assert.Equal(t, true, opt2.isWrapped())
}

func TestUnwrap(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	opt.Exists(transform)

	opt2 := Wrap(opt)

	// Act
	innerOpt := Unwrap(opt2)

	// Assert
	assert.Equal(t, "hello", innerOpt.Greeting)
}

func TestDeeplyWrappedUnwrap(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	opt.Exists(transform)

	opt2 := Wrap(opt)

	opt3 := Wrap(opt2)

	// Act
	innerOpt := Unwrap(opt3)

	// Assert
	assert.Equal(t, "hello", innerOpt.Greeting)
}
