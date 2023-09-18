# goptional
[![Go Reference](https://pkg.go.dev/badge/github.com/syke99/goptional.svg)](https://pkg.go.dev/github.com/syke99/goptional)
[![go reportcard](https://goreportcard.com/badge/github.com/syke99/goptional)](https://goreportcard.com/report/github.com/syke99/goptional)
[![License](https://img.shields.io/github/license/syke99/goptional)](https://github.com/syke99/goptional/blob/master/LICENSE)
![Go version](https://img.shields.io/github/go-mod/go-version/syke99/goptional)</br>
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
    opt.Map(transform)
    
    // Assert
    assert.Equal(t, "hello", tt.greeting) // asserts true
}
```

Who?
====

This library was developed by Quinn Millican ([@syke99](https://github.com/syke99))

