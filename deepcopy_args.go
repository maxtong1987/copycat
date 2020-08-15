package copycat

import (
	"reflect"
)

type visitedAddr struct {
	a uintptr
	t reflect.Type
}

type deepCopyArgs struct {
	d        reflect.Value
	s        reflect.Value
	flags    Flags
	level    uint
	maxLevel uint
	visited  *map[visitedAddr]reflect.Value
}

func (args *deepCopyArgs) resolve() *deepCopyArgs {
	flags := args.flags
	args.d = resolveDst(args.d, flags)
	args.s = resolveSrc(args.s, flags)
	return args
}

func (args *deepCopyArgs) next() *deepCopyArgs {
	nextArgs := *args
	nextArgs.level++
	return &nextArgs
}

func (args *deepCopyArgs) recordVisited(addr visitedAddr) {
	(*args.visited)[addr] = args.d
}

func resolveDst(v reflect.Value, flags Flags) reflect.Value {
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

func resolveSrc(v reflect.Value, flags Flags) reflect.Value {
	for {
		switch k := v.Kind(); k {
		case reflect.Ptr, reflect.Interface:
			v = v.Elem()
		default:
			return v
		}
	}
}
