package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd is the main "ghpr" command
var rootCmd = &cobra.Command{
	Use:   "ghpr",
	Short: "create and search github PRs and issues from the command line",
}

func init() {
	cfgFile := fmt.Sprintf("%s/.ghpr.json", os.Getenv("HOME"))
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	Token = viper.GetString("token")
	inEditor = viper.GetBool("inEditor")

}

// Execute executes the command and returns the error
func Execute() error {
	return rootCmd.Execute()
}
