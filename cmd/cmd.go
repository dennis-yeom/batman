package cmd

import (
	"fmt"
	"log"

	"github.com/dennis-yeom/batman/internal/demo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	port int

	// RootCmd is the main command for the CLI
	RootCmd = &cobra.Command{
		Use:   "demo",
		Short: "start the demo",
		Long:  "Starts the demo execution for Dennis",
		Run: func(cmd *cobra.Command, args []string) {
			//fmt.Println("Running the demo")
		},
	}

	// test connection to redis
	TestCmd = &cobra.Command{
		Use:   "test",
		Short: "tests connection to redis",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := demo.New(port)

			if err != nil {
				return err
			}

			// Proceed with using d since it was created successfully
			fmt.Println("Demo instance created successfully:", d)
			return nil
		},
	}
)

func init() {
	// Add the TestCmd to the RootCmd
	RootCmd.AddCommand(TestCmd)

	// Set up Viper to read configuration from .config.yml
	viper.SetConfigName(".config") // name of config file (without extension)
	viper.SetConfigType("yaml")    // required since we're using .yml
	viper.AddConfigPath(".")       // look for the config file in the current directory

	// Set default values in case .config.yml does not exist or lacks specific entries
	viper.SetDefault("redis.port", 6380)

	// Bind Viper values to flags
	RootCmd.PersistentFlags().IntVarP(&port, "port", "p", viper.GetInt("redis.port"), "port of redis cache")

	// Load the config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No configuration file found; using defaults or command-line args: %v", err)
	}

	// Bind Viper keys to flags so changes reflect in CLI options
	viper.BindPFlags(RootCmd.PersistentFlags())
}

// Execute runs the RootCmd and handles any errors
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
