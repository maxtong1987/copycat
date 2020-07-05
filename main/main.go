package main

import (
	"fmt"
	"reflect"
)

type SA struct {
	A string
	B int
	C float32
}

func (s SA) Foo() {
	fmt.Println("this is foo")
}

type SB struct {
	SA
	D uint
}

type Level uint

func (l Level) Next() Level { return Level(l + 1) }

func (l Level) Indent() string {
	bytes := make([]byte, l)
	for i := 0; i < int(l); i++ {
		bytes[i] = 9 // tab
	}
	return string(bytes)
}

func walker(in interface{}) {
	_walker(reflect.ValueOf(in), 0)
}

func _walker(value reflect.Value, level Level) {
	if !value.IsValid() {
		return
	}
	k := value.Kind()
	t := value.Type()
	fmt.Printf("%sk: %s  t: %s\n", level.Indent(), k, t)
	switch k {
	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			_walker(value.Index(i), level.Next())
		}
	case reflect.Interface, reflect.Ptr:
		_walker(value.Elem(), level.Next())
	case reflect.Struct:
		numField := value.NumField()
		for i := 0; i < numField; i++ {
			f := t.Field(i)
			fmt.Printf("%s%s:\n", level.Indent(), f.Name)
			_walker(value.Field(i), level.Next())
		}
		numMethod := value.NumMethod()
		for i := 0; i < numMethod; i++ {
			m := t.Method(i)
			fmt.Printf("%s%s:\n", level.Indent(), m.Name)
			_walker(value.Method(i), level.Next())
		}
	case reflect.Map:
		for iter := value.MapRange(); iter.Next(); {
			key := reflect.New(t.Key()).Elem()
			fmt.Printf("%s%key:\n", level.Indent())
			_walker(key, level.Next())
			value := reflect.New(t.Elem()).Elem()
			fmt.Printf("%s%value:\n", level.Indent())
			_walker(value, level.Next())
		}
	default:
		return
	}
}

func main() {
	sb := SB{
		SA: SA{
			A: "a",
			B: 123,
			C: 123.3,
		},
		D: 444,
	}
	walker(sb)
}
