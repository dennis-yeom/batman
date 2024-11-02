package cmd

import (
	"fmt"
	"log" //import go library for logging. used for logging errors.

	"github.com/dennis-yeom/batman/internal/demo" //import

	"github.com/spf13/cobra" //cobra & viper libraries used to create CLI
	"github.com/spf13/viper"
)

var (
	port  int
	key   string
	value string

	// RootCmd is the main command for the CLI
	RootCmd = &cobra.Command{
		Use:   "demo",
		Short: "start the demo",
		Long:  "Starts the demo execution for Dennis",
		Run: func(cmd *cobra.Command, args []string) { //run this command if args is empty
			fmt.Println("Running the demo. \nFor help: go run main.go --help")
		},
	}

	// SetCmd sets a key and value in Redis
	SetCmd = &cobra.Command{ //create a new cobra command names SetCmd
		Use:   "set",                //defines name if command to be used in CLI: go run main.go set
		Short: "sets key and value", //description of the commands function, appears when using --help
		RunE: func(cmd *cobra.Command, args []string) error { //when 'set' is used, run this with error handling
			d, err := demo.New(port) //create new demo instance using port number.
			if err != nil {          //checking if we successfully created the instance
				return err
			}
			return d.Set(key, value) //run the set function with provided key, value and return result. defined in demo
		},
	}

	// GetCmd retrieves a value from Redis based on the key
	GetCmd = &cobra.Command{
		Use:   "get",
		Short: "gets value for key",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := demo.New(port)
			if err != nil {
				return err
			}
			return d.Get(key)
		},
	}
)

// set up viper for easy configuration
func init() {
	viper.SetConfigName(".config") // name of config file (without extension)
	viper.SetConfigType("yaml")    // required since we're using .yml
	viper.AddConfigPath(".")       // look for the config file in the current directory

	// Set default values in case .config.yml does not exist or lacks specific entries
	viper.SetDefault("redis.port", 6380)

	// Load the config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No configuration file found; using defaults or command-line args: %v", err)
	}

	// Add the TestCmd to the RootCmd
	RootCmd.AddCommand(SetCmd)
	RootCmd.AddCommand(GetCmd)

	// Bind Viper values to flags
	RootCmd.PersistentFlags().IntVarP(&port, "port", "p", viper.GetInt("redis.port"), "port of redis cache")

	// Flags for SetCmd
	SetCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "name of the key")
	SetCmd.PersistentFlags().StringVarP(&value, "value", "v", "", "name of the value")

	// Flags for GetCmd
	GetCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "name of the key")

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
