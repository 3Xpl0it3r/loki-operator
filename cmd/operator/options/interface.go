package options

import "github.com/spf13/pflag"

// all custom options should implement this interfaces
type options interface {
	Validate() []error
	Complete() error
	AddFlags(*pflag.FlagSet)
}
