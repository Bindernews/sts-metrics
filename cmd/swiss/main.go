package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bindernews/sts-msr/tools"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
	if err := rootCommand(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func rootCommand(args []string) error {
	commands := []tools.ICommand{
		tools.NewArchiveExportCmd(),
	}
	// Make sure we have at least one arg, so we can get through
	// the loop and print the subcommand names
	if len(args) < 1 {
		args = append(args, "")
	}
	subcommand := args[0]
	names := make([]string, 0)
	for _, cmd := range commands {
		cmdName := cmd.Flags().Name()
		names = append(names, cmdName)
		if cmdName == subcommand {
			if err := cmd.Flags().Parse(args[1:]); err != nil {
				return err
			}
			if err := cmd.Run(); err != nil {
				return err
			}
			return nil
		}
	}
	// No command found
	return fmt.Errorf("subcommands: %s", strings.Join(names, ", "))
}
