package goptional

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func transform(val *testType) {
	val.Greeting = "hi all"
}

func fail(val *testType) {
	val.fail = true
}

func checkFailure(val *testType) error {
	if val.fail {
		return errors.New("failed successfully")
	}
	return nil
}

type testType struct {
	Greeting string `json:"greeting"`
	fail     bool
}

func TestNewGoptional(t *testing.T) {
	// Arrange
	tt := testType{}

	// Act
	opt := NewGoptional(&tt)

	// Assert
	assert.NotNil(t, opt)
}

func TestWrapUnwrap(t *testing.T) {
	// Arrange
	tt := testType{Greeting: "hello"}

	opt := NewGoptional(&tt)

	// Act
	opt = Wrap(opt)

	ptr := Unwrap(opt)

	// Assert
	assert.Equal(t, "hello", ptr.Greeting)
}

func TestMap(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	// Act
	opt.Map(transform)

	// Assert
	assert.Equal(t, "hi all", tt.Greeting)
}

func TestMapElseExists(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	// Act
	opt.MapElse(transform, func() *testType {
		return &testType{}
	})

	// Assert
	assert.Equal(t, "hi all", tt.Greeting)
}

func TestMapElseDoesntExists(t *testing.T) {
	// Arrange
	opt := NewGoptional[testType](nil)

	// Act
	opt.MapElse(transform, func() *testType {
		return &testType{}
	})

	// Assert
	assert.Equal(t, "hi all", any(opt.Val()).(testType).Greeting)
}

func TestMapElseDoesntExistString(t *testing.T) {
	// Arrange
	greeting := ""

	opt := NewGoptional(&greeting)

	opt.ExistsNil()

	opt.MapElse(func(s *string) {
		if *s == "hi" {
			*s = "hello"
		}
	}, func() *string {
		greet := "hi"
		return &greet
	})
}

func TestExistsNil(t *testing.T) {
	// Arrange
	greeting := ""

	opt := NewGoptional(&greeting)

	// Act
	opt.ExistsNil()

	// Assert
	assert.Nil(t, opt.unwrapVal())
}

func TestMapElseDoesntExistAfterNil(t *testing.T) {
	// Arrange
	tt := testType{Greeting: "hello"}

	opt := NewGoptional(&tt)

	// Act
	opt.ExistsNil()

	opt.MapElse(transform, func() *testType {
		return &testType{}
	})

	// Assert
	assert.Equal(t, "hi all", any(Unwrap(opt)).(testType).Greeting)
}

func TestExistsDoesExist(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	// Act
	err := opt.Exists(checkFailure)

	// Assert
	assert.NoError(t, err)
}

func TestExistsDoesntExist(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	// Act
	err := opt.Map(fail).Exists(checkFailure)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "failed successfully", err.Error())
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

	opt := NewGoptional[testType](nil)

	// Act
	v := opt.ValOr(tt)

	// Assert
	assert.Equal(t, "hello", any(v).(testType).Greeting)
}

func TestValElse(t *testing.T) {
	// Arrange
	opt := NewGoptional[testType](nil)

	// Act
	v := opt.ValElse(func() testType {
		return testType{Greeting: "hello"}
	})

	// Assert
	assert.Equal(t, "hello", v.(testType).Greeting)
}

func TestMarshalJSON(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	opt.Map(transform)

	// Act
	jsonBytes, err := opt.MarshalJSON()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "{\"greeting\":\"hi all\"}", string(jsonBytes))
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

func TestUnmarshalJSON(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	jsonString := "{\"greeting\":\"hello\"}"

	// Act
	err := opt.UnmarshalJSON([]byte(jsonString))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "hello", any(opt.Val()).(testType).Greeting)
}

func TestUnmarshalJSONDoesntExist(t *testing.T) {
	// Arrange
	opt := NewGoptional[*testType](nil)

	jsonString := "{\"greeting\":\"hello\"}"

	// Act
	err := opt.UnmarshalJSON([]byte(jsonString))

	// Assert
	assert.NoError(t, err)
}

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

	opt.Map(transform)

	opt2 := Wrap(opt)

	// Act
	innerOpt := Unwrap(opt2)

	// Assert
	assert.Equal(t, "hi all", innerOpt.Greeting)
}

func TestDeeplyWrappedUnwrap(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	opt.Map(transform)

	opt2 := Wrap(opt)

	opt3 := Wrap(opt2)

	// Act
	innerOpt := Unwrap(opt3)

	// Assert
	assert.Equal(t, "hi all", innerOpt.Greeting)
}

func TestUnmarshalJSONUnwrap(t *testing.T) {
	// Arrange
	tt := testType{}

	opt := NewGoptional(&tt)

	opt = Wrap(opt)

	jsonString := "{\"greeting\":\"hello\"}"

	// Act
	err := opt.UnmarshalJSON([]byte(jsonString))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "hello", tt.Greeting)
}
