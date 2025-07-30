package cli

import "strings"

// SetValues is a custom flag type for handling multiple --set flags.
type SetValues []string

func (s *SetValues) String() string {
	return strings.Join(*s, ",")
}

func (s *SetValues) Set(value string) error {
	*s = append(*s, value)
	return nil
}
