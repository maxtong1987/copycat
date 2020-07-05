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
	"reflect"
)

// Flags Bitmask for DeepCopy
type Flags uint32

const (
	// FCopyChan true if copy from source's channel to destination's channel
	FCopyChan Flags = 1 << iota
	// FMapToStruct true if copy source's map to destination's struct
	FMapToStruct
)

// Has check if f has flags
func (f Flags) Has(in Flags) bool { return f&in != 0 }

// DeepCopy recursively copies data from src to dst.
// Doesn't support: cannel, function and unsafe pointer
func DeepCopy(dst interface{}, src interface{}, flags ...Flags) {
	var flagsCombined Flags
	for _, f := range flags {
		flagsCombined |= f
	}
	deepCopy(reflect.ValueOf(dst), reflect.ValueOf(src), flagsCombined)
}

func deepCopy(dst reflect.Value, src reflect.Value, flags Flags) {

	d := indirect(dst, true)
	s := indirect(src, false)

	if !canCopy(d, s) {
		return
	}

	switch d.Kind() {

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
		structHandler(d, s, flags)

	case reflect.Map:
		mapHandler(d, s, flags)

	case reflect.Array:
		arrayHandler(d, s, flags)

	case reflect.Slice:
		sliceHandler(d, s, flags)

	case reflect.Chan:
		d.Set(s)

	default: // Func, UnsafePointer ...
		return
	}
}

func indirect(v reflect.Value, isDst bool) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if isDst && v.IsNil() {
			ptr := reflect.New(v.Type().Elem())
			v.Set(ptr)
		}
		return indirect(v.Elem(), isDst)
	default:
		return v
	}
}

func structHandler(d reflect.Value, s reflect.Value, flags Flags) {
	t := d.Type()
	num := t.NumField()
	for i := 0; i < num; i++ {
		f := t.Field(i)
		df := d.Field(i)
		sf := s.FieldByName(f.Name)
		deepCopy(df, sf, flags)
	}
}

func mapHandler(d reflect.Value, s reflect.Value, flags Flags) {
	t := d.Type()
	newMap := reflect.MakeMap(t)
	d.Set(newMap)
	for iter := s.MapRange(); iter.Next(); {
		key := reflect.New(t.Key()).Elem()
		value := reflect.New(t.Elem()).Elem()
		deepCopy(key, iter.Key(), flags)
		deepCopy(value, iter.Value(), flags)
		d.SetMapIndex(key, value)
	}
}

func sliceHandler(d reflect.Value, s reflect.Value, flags Flags) {
	t := d.Type()
	arr := reflect.MakeSlice(t, s.Len(), s.Cap())
	d.Set(arr)
	copyArr(d, s, flags)
}

func arrayHandler(d reflect.Value, s reflect.Value, flags Flags) {
	copyArr(d, s, flags)
}

func copyArr(d reflect.Value, s reflect.Value, flags Flags) {
	len := d.Len()
	if len > s.Len() {
		len = s.Len()
	}
	for i := 0; i < len; i++ {
		deepCopy(d.Index(i), s.Index(i), flags)
	}
}

func isInt(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	default:
		return false
	}
}

func isUint(k reflect.Kind) bool {
	switch k {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func isComplex(k reflect.Kind) bool {
	switch k {
	case reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}

func isFloat(k reflect.Kind) bool {
	switch k {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func isArrayOrSlice(k reflect.Kind) bool {
	switch k {
	case reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}

func canCopy(d reflect.Value, s reflect.Value) bool {
	if !d.IsValid() || !s.IsValid() || !d.CanSet() {
		return false
	}
	dk := d.Kind()
	sk := s.Kind()
	if d.Kind() == s.Kind() {
		return true
	}
	if isInt(dk) && isInt(sk) {
		return true
	}
	if isUint(dk) && isUint(sk) {
		return true
	}
	if isFloat(dk) && isFloat(sk) {
		return true
	}
	if isComplex(dk) && isComplex(sk) {
		return true
	}
	if isArrayOrSlice(dk) && isArrayOrSlice(sk) {
		return true
	}
	return false
}
