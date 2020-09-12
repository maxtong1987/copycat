# copycat
[![Coverage Status](https://coveralls.io/repos/github/maxtong1987/copycat/badge.svg?branch=master)](https://coveralls.io/github/maxtong1987/copycat?branch=master)

Recursively perform deep copy from source (src) to destination (dst) using reflection.

## Getting started

To get the package, execute:
```sh
go get github.com/maxtong1987/copycat
```

To import the package, add the following line to your code:
```go
import "github.com/maxtong1987/copycat"
```

## Example
```go
package main

import (
	"fmt"

	"github.com/maxtong1987/copycat"
)

type subType struct {
	X string
	Y bool
	Z int
}

type srcType struct {
	A string
	B int32
	C float64
	D []uint64
	E subType
}

type destType struct {
	A string
	B int64
	C float32
	d []uint64
	E *subType
}

func (d *destType) String() string {
	return fmt.Sprintf("A:%s B:%v C:%v d:%v E:%+v", d.A, d.B, d.C, d.d, *d.E)
}

func main() {
	src := srcType{
		A: "a",
		B: 123,
		C: 0.0122,
		D: []uint64{6, 7, 8, 9},
		E: subType{X: "x", Y: true, Z: 100},
	}
	dst := &destType{}
	copycat.DeepCopy(dst, src)
	fmt.Print(dst)
}
```

## Expected Behavior

|Source|Destination|Expected|
|---|---|---|
|struct{A:"a",B:123}|struct{}|struct{A:"a",B:123}|
|struct{A:"a",B:123}|*struct{}|*struct{A:"a",B:123}|
|struct{A:"a",B:123}|nil|*struct{A:"a",B:123}|
|**struct{A:"a",B:123}|nil|*struct{A:"a",B:123}|
|struct{a:"a",b:123}|struct{}|struct{}|
|nil|nil|nil|
