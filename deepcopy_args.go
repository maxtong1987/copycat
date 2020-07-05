package copycat

import (
	"reflect"
)

type deepCopyArgs struct {
	d       reflect.Value
	s       reflect.Value
	flags   Flags
	level   uint
	visited *map[uintptr]reflect.Value
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

func (args *deepCopyArgs) recordVisited() {
	s := args.s
	d := args.d
	flags := args.flags
	visited := *args.visited
	if !flags.Has(FPreserveHierarchy) || !s.CanAddr() {
		return
	}
	addr := s.UnsafeAddr()
	visited[addr] = d
}

func resolveDst(v reflect.Value, flags Flags) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsValid() && v.IsNil() {
			ptr := reflect.New(v.Type().Elem())
			v.Set(ptr)
		}
		return resolveDst(v.Elem(), flags)
	case reflect.Interface:
		// if v.IsValid() && v.IsNil() {
		// 	ptr := reflect.New(v.Type().Elem())
		// 	v.Set(ptr)
		// }
		return resolveDst(v.Elem(), flags)
	default:
		return v
	}
}

func resolveSrc(v reflect.Value, flags Flags) reflect.Value {
	switch k := v.Kind(); k {
	case reflect.Ptr:
		return resolveSrc(v.Elem(), flags)
	case reflect.Interface:
		return resolveSrc(v.Elem(), flags)
	default:
		return v
	}
}
