package copycat

// Flags Bitmask for DeepCopy
type Flags uint32

const (

	// FPreserveHierarchy preserve pointer hierarchy
	FPreserveHierarchy Flags = 1 << iota

	// FCopyChan copy from source's channel to destination's channel
	FCopyChan

	// FCopyPtr copy pointer
	FCopyPtr

	// FCopyFunc copy function
	FCopyFunc

	//reflect.Ptr, reflect.Uintptr, reflect.Chan, reflect.Func, reflect.UnsafePointer:

	// FCopyUintptr copy Uintptr
	FCopyUintptr

	// FCopyUnsafePointer copy UnsafePointer
	FCopyUnsafePointer
)

// Has check if f has flags
func (f Flags) Has(in Flags) bool { return f&in != 0 }

// Combine combines all flags mask
func combineFlags(flags ...Flags) Flags {
	var combined Flags
	for _, f := range flags {
		combined |= f
	}
	return combined
}
