package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	// RootCmd is the main command for the CLI
	RootCmd = &cobra.Command{
		Use:   "demo",
		Short: "start the demo",
		Long:  "Starts the demo execution for Dennis",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Running the demo")
		},
	}

	// create a teset command
	TestCmd = &cobra.Command{
		Use:   "test",
		Short: "this is the command i made",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("You chose the test command!")
			return nil
		},
	}
)

func init() {
	// Add the TestCmd to the RootCmd
	RootCmd.AddCommand(TestCmd)
}

// Execute runs the RootCmd and handles any errors
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
