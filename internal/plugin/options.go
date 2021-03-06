package plugin

import (
	"fmt"
	"strings"
)

type Options struct {
	// The filename the generator will output into.
	// Defaults to `"index.ts"`.
	Filename string

	// The insertion point in the target file the generator
	// use.
	// Defaults to `""` (ie. no insertion point)
	InsertionPoint string
}

func defaultOptions() Options {
	return Options{
		Filename:       "index.ts",
		InsertionPoint: "",
	}
}

func (o *Options) Unmarshal(s *string) error {
	defaults := defaultOptions()
	o.Filename = defaults.Filename
	o.InsertionPoint = defaults.InsertionPoint

	// no options specified
	if s == nil {
		return nil
	}
	str := *s

	opts := strings.Split(str, ",")
	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid option [%s]: expected key and value", opt)
		}
		key, value := parts[0], parts[1]
		switch key {
		case "insertion_point":
			o.InsertionPoint = value
		case "filename":
			o.Filename = value
		default:
			return fmt.Errorf("unknown option [%s]", opt)
		}
	}
	return nil
}
