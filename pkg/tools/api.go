package tools

import "flag"

type ICommand interface {
	Flags() *flag.FlagSet
	// Returns the usage description for the command
	Description() string
	// Run the command
	Run() error
}
