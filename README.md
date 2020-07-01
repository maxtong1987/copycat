# copycat

Recursively perform deep copy from source (src) to destination (dst) using reflection.

## Getting started

To get the package, execute:
```sh
go get gitlab.com/maxtong/copycat
```

To import the package, add the following line to your code:
```go
import "gitlab.com/maxtong/copycat"
```

## Example
```go
package main

import (
    "fmt"
    "gitlab.com/maxtong/copycat"
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

func main() {
    src := srcType{
        A: "a",
        B: 123,
        C: 0.0122,
        D: []uint64{6,7,8,9}
        E: subType{X: "x", Y: true, Z: 100}
    }
    dst := &destType{}
    copycat.DeepCopy(dst, src)
    fmt.Printf("dst:%+v", *dst)
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
