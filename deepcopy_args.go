package copycat

import "reflect"

type deepCopyArgs struct {
	d     reflect.Value
	s     reflect.Value
	flags Flags
	level uint
}

func (args *deepCopyArgs) indirect() *deepCopyArgs {
	args.d = indirect(args.d, true)
	args.s = indirect(args.s, false)
	return args
}

func (args *deepCopyArgs) next() *deepCopyArgs {
	nextArgs := *args
	nextArgs.level++
	return &nextArgs
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
