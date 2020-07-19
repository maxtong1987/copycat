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
func DeepCopy(dst interface{}, src interface{}, flags ...Flags) error {
	args := deepCopyArgs{
		d:       reflect.ValueOf(dst),
		s:       reflect.ValueOf(src),
		flags:   combineFlags(flags...),
		level:   0,
		visited: &map[uintptr]reflect.Value{},
	}
	return deepCopy(&args)
}

func deepCopy(args *deepCopyArgs) error {

	args.resolve()
	d := args.d
	s := args.s
	flags := args.flags

	if !canCopy(d, s) {
		return nil
	}

	if flags.Has(FPreserveHierarchy) {
		if s.CanAddr() {
			addr := s.UnsafeAddr()
			if value, ok := (*args.visited)[addr]; ok {
				d.Set(value)
				return nil
			}
			defer args.recordVisited(addr)
		}
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
		return structHandler(args)

	case reflect.Map:
		return mapHandler(args)

	case reflect.Array:
		return arrayHandler(args)

	case reflect.Slice:
		return sliceHandler(args)

	case reflect.Chan:
		if flags.Has(FCopyChan) {
			d.Set(s)
		}

	case reflect.Func:
		if flags.Has(FCopyFunc) {
			d.Set(s)
		}

	case reflect.Uintptr:
		if flags.Has(FCopyUintptr) {
			d.Set(s)
		}

	case reflect.UnsafePointer:
		if flags.Has(FCopyUnsafePointer) {
			d.Set(s)
		}

	case reflect.Interface:
		if flags.Has(FCopyInterface) {
			d.Set(s)
		}

	default: // Invalid
		return fmt.Errorf("unhandled type: %s", k)
	}

	return nil
}

func structHandler(args *deepCopyArgs) error {
	d := args.d
	s := args.s
	t := d.Type()
	for i, num := 0, t.NumField(); i < num; i++ {
		f := t.Field(i)
		nextArgs := args.next()
		nextArgs.d = d.Field(i)
		nextArgs.s = s.FieldByName(f.Name)
		if err := deepCopy(nextArgs); err != nil {
			return err
		}
	}
	for i, num := 0, t.NumMethod(); i < num; i++ {
		m := t.Method(i)
		nextArgs := args.next()
		nextArgs.d = d.Method(i)
		nextArgs.s = s.MethodByName(m.Name)
		if err := deepCopy(nextArgs); err != nil {
			return err
		}
	}
	return nil
}

func mapHandler(args *deepCopyArgs) error {
	d := args.d
	s := args.s
	t := d.Type()
	newMap := reflect.MakeMap(t)
	d.Set(newMap)
	for iter := s.MapRange(); iter.Next(); {
		nextArgs := args.next()
		nextArgs.d = reflect.New(t.Elem()).Elem()
		nextArgs.s = iter.Value()
		if err := deepCopy(nextArgs); err != nil {
			return err
		}
		d.SetMapIndex(iter.Key(), nextArgs.d)
	}
	return nil
}

func sliceHandler(args *deepCopyArgs) error {
	d := args.d
	s := args.s
	t := d.Type()
	len := s.Len()
	arr := reflect.MakeSlice(t, len, s.Cap())
	d.Set(arr)
	dk := d.Type().Elem().Kind()
	sk := s.Type().Elem().Kind()
	if dk == reflect.Uint8 && sk == reflect.Uint8 {
		d.SetBytes(s.Bytes())
		return nil
	}
	for i := 0; i < len; i++ {
		nextArgs := args.next()
		nextArgs.d = d.Index(i)
		nextArgs.s = s.Index(i)
		if err := deepCopy(nextArgs); err != nil {
			return err
		}
	}
	return nil
}

func arrayHandler(args *deepCopyArgs) error {
	d := args.d
	s := args.s
	len := d.Len()
	if len > s.Len() {
		len = s.Len()
	}
	for i := 0; i < len; i++ {
		nextArgs := args.next()
		nextArgs.d = d.Index(i)
		nextArgs.s = s.Index(i)
		if err := deepCopy(nextArgs); err != nil {
			return err
		}
	}
	return nil
}
