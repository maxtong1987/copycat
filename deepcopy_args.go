package copycat

import (
	"reflect"
)

type visitedAddr struct {
	a uintptr
	t reflect.Type
}

type deepCopyArgs struct {
	d       reflect.Value
	s       reflect.Value
	visited *map[visitedAddr]reflect.Value
}

func (args *deepCopyArgs) resolve() *deepCopyArgs {
	args.d = resolveDst(args.d)
	args.s = resolveSrc(args.s)
	return args
}

func (args *deepCopyArgs) next() *deepCopyArgs {
	nextArgs := *args
	return &nextArgs
}

func (args *deepCopyArgs) recordVisited(addr visitedAddr) {
	(*args.visited)[addr] = args.d
}

func resolveDst(v reflect.Value) reflect.Value {
	for {
		if v.Kind() != reflect.Ptr {
			return v
		}
		if v.IsNil() && v.CanSet() {
			ptr := reflect.New(v.Type().Elem())
			v.Set(ptr)
		}
		v = v.Elem()
	}
}

func resolveSrc(v reflect.Value) reflect.Value {
	for {
		switch k := v.Kind(); k {
		case reflect.Ptr, reflect.Interface:
			v = v.Elem()
		default:
			return v
		}
	}
}
