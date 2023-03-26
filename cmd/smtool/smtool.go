package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bindernews/sts-msr/pkg/tools"
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
		tools.NewUploadRunsCmd(),
	}
	// Make sure we have at least one arg, so we can get through
	// the loop and print the subcommand names
	if len(args) < 1 {
		args = append(args, "")
	}
	subcommand := args[0]
	for _, cmd := range commands {
		cmdName := cmd.Flags().Name()
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
	if args[0] == "-h" || args[0] == "--help" {
		fmt.Println("subcommands:")
		for _, cmd := range commands {
			fmt.Printf("  %-12s %s\n", cmd.Flags().Name(), cmd.Description())
		}
		return nil
	} else {
		return fmt.Errorf("use --help for usage")
	}
}
