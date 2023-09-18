# goptional
an Option package for Go

How do I use goptional?
====

### Installation

```
go get github.com/DragonsDenSoftware/goptional
```

### Basic Usage

```go
package main

import (
	"testing"
    
	"github.com/syke99/goptional"
	"github.com/stretchr/testify/assert"
)

func transform(val *testType) {
	val.greeting = "hello"
}

type testType struct {
	greeting string
}

func main() {
    // Arrange
    tt := testType{}
    
    opt := NewGoptional(&tt)
    
    // Act
    opt.Exists(transform)
    
    // Assert
    assert.Equal(t, "hello", tt.greeting) // asserts true
}
```

Who?
====

This library was developed by Quinn Millican ([@syke99](https://github.com/syke99))

