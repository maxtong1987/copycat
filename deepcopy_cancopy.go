package copycat

import "reflect"

func canCopy(args *deepCopyArgs) bool {
	if args == nil {
		return false
	}

	d := args.d
	s := args.s
	flags := args.flags

	if !d.IsValid() || !s.IsValid() || !d.CanSet() {
		return false
	}

	dk := d.Kind()
	sk := s.Kind()
	if dk == sk {
		switch dk {
		case reflect.Chan:
			return flags.Has(FCopyChan)
		case reflect.Func:
			return flags.Has(FCopyFunc)
		case reflect.Ptr:
			return flags.Has(FCopyPtr)
		case reflect.Uintptr:
			return flags.Has(FCopyUintptr)
		case reflect.UnsafePointer:
			return flags.Has(FCopyUnsafePointer)
		default:
			return true
		}
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
