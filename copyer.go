package copycat

import "reflect"

// Copyer a configurable deepcopy helper
type Copyer struct {
	flags    Flags
	maxLevel uint
}

// GetFlags returns flags
func (c *Copyer) GetFlags() Flags { return c.flags }

// SetFlags set flags
func (c *Copyer) SetFlags(flags Flags) *Copyer {
	c.flags = flags
	return c
}

// GetMaxlevel returns maximum copy level
func (c *Copyer) GetMaxlevel() uint { return c.maxLevel }

// SetMaxLevel set maximum copy level
func (c *Copyer) SetMaxLevel(maxLevel uint) *Copyer {
	c.maxLevel = maxLevel
	return c
}

// NewCopyer returns Coper pointer with default values
func NewCopyer() *Copyer {
	return &Copyer{
		flags:    0,
		maxLevel: DefaultMaxLevel,
	}
}

// DeepCopy deep copies data from src to dst
func (c *Copyer) DeepCopy(dst interface{}, src interface{}) error {
	args := deepCopyArgs{
		d:        reflect.ValueOf(dst),
		s:        reflect.ValueOf(src),
		flags:    c.flags,
		level:    0,
		maxLevel: c.maxLevel,
		visited:  &map[visitedAddr]reflect.Value{},
	}
	return deepCopy(&args)
}
