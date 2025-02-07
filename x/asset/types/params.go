package types

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		AllowExtensions: []string{"transfer", "mint"},
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params
func (p Params) Validate() error {
	for _, extension := range p.AllowExtensions {
		if extension == "" {
			return fmt.Errorf("Extension can not be empty")
		}
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
