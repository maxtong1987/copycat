package copycat

// Flags Bitmask for DeepCopy
type Flags uint16

const (

	// FPreserveHierarchy preserve pointer hierarchy
	FPreserveHierarchy Flags = 1 << iota

	// FCopyChan shallow copy Chan type
	FCopyChan

	// FCopyFunc shallow copy Func type
	FCopyFunc

	// FCopyUintptr shallow copy Uintptr type
	FCopyUintptr

	// FCopyUnsafePointer shallow copy UnsafePointer type
	FCopyUnsafePointer

	// FCopyInterface shallow copy interface type
	FCopyInterface

	// FAll shallow copy all special types
	FAll = ^Flags(0)
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
