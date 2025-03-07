package types

import (
	"gopkg.in/yaml.v2"
)

// NewParams returns Params instance with the given values.
func NewParams(authority string) Params {
	return Params{
		Authority: authority,
	}
}

// default bridge module parameters
func DefaultParams() Params {
	return Params{
		Authority: "",
	}
}

// validate params
func (p Params) Validate() error {
	// TODO: validate
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
