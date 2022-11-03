package flag

import (
	"fmt"
	"strings"
)

type enumFlag struct {
	Allowed []string
	Value   string
}

// Provides the logic for handling a flag that can only assume a value in a defined set of values
// See https://github.com/spf13/pflag/issues/236
func NewEnumFlag(allowed []string, defaultValue string) *enumFlag {
	return &enumFlag{
		Allowed: allowed,
		Value:   defaultValue,
	}
}

func (a enumFlag) String() string {
	return a.Value
}

func (a *enumFlag) Set(p string) error {
	isIncluded := func(opts []string, val string) bool {
		for _, opt := range opts {
			if val == opt {
				return true
			}
		}
		return false
	}
	if !isIncluded(a.Allowed, p) {
		return fmt.Errorf("%s is not included in %s", p, strings.Join(a.Allowed, ","))
	}
	a.Value = p
	return nil
}

func (a *enumFlag) Type() string {
	return "string"
}
