package tools

import "flag"

type ICommand interface {
	Flags() *flag.FlagSet
	Run() error
}
