/**
Copyright (c) 2020 Max Tong

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
**/

package copycat

import (
	"fmt"
	"reflect"
)

// DeepCopy recursively copies data from src to dst.
// Doesn't support: cannel, function and unsafe pointer
func DeepCopy(dst interface{}, src interface{}, flags ...Flags) {
	args := deepCopyArgs{
		d:     reflect.ValueOf(dst),
		s:     reflect.ValueOf(src),
		flags: combineFlags(flags...),
		level: 0,
	}
	deepCopy(&args)
}

func deepCopy(args *deepCopyArgs) {

	args.indirect()
	d := args.d
	s := args.s

	if !canCopy(args) {
		return
	}

	switch k := d.Kind(); k {

	case reflect.String:
		d.SetString(s.String())

	case reflect.Bool:
		d.SetBool(s.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.SetInt(s.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d.SetUint(s.Uint())

	case reflect.Float32, reflect.Float64:
		d.SetFloat(s.Float())

	case reflect.Complex64, reflect.Complex128:
		d.SetComplex(s.Complex())

	case reflect.Struct:
		structHandler(args)

	case reflect.Map:
		mapHandler(args)

	case reflect.Array:
		arrayHandler(args)

	case reflect.Slice:
		sliceHandler(args)

	case reflect.Ptr, reflect.Uintptr, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		d.Set(s)

	default:
		panic(fmt.Sprintf("unhandled type: %s", k))
	}
}

func structHandler(args *deepCopyArgs) {
	d := args.d
	s := args.s
	t := d.Type()
	num := t.NumField()
	for i := 0; i < num; i++ {
		f := t.Field(i)
		nextArgs := args.next()
		nextArgs.d = d.Field(i)
		nextArgs.s = s.FieldByName(f.Name)
		deepCopy(nextArgs)
	}
}

func mapHandler(args *deepCopyArgs) {
	d := args.d
	s := args.s
	t := d.Type()
	newMap := reflect.MakeMap(t)
	d.Set(newMap)
	for iter := s.MapRange(); iter.Next(); {
		nextArgs := args.next()
		nextArgs.d = reflect.New(t.Elem()).Elem()
		nextArgs.s = iter.Value()
		deepCopy(nextArgs)
		d.SetMapIndex(iter.Key(), nextArgs.d)
	}
}

func sliceHandler(args *deepCopyArgs) {
	d := args.d
	s := args.s
	t := d.Type()
	arr := reflect.MakeSlice(t, s.Len(), s.Cap())
	d.Set(arr)
	copyArr(args)
}

func arrayHandler(args *deepCopyArgs) {
	copyArr(args)
}

func copyArr(args *deepCopyArgs) {
	d := args.d
	s := args.s
	len := d.Len()
	if len > s.Len() {
		len = s.Len()
	}
	dk := d.Type().Elem().Kind()
	sk := s.Type().Elem().Kind()
	if dk == reflect.Uint8 && sk == reflect.Uint8 {
		d.SetBytes(s.Bytes())
		return
	}
	for i := 0; i < len; i++ {
		nextArgs := args.next()
		nextArgs.d = d.Index(i)
		nextArgs.s = s.Index(i)
		deepCopy(nextArgs)
	}
}
