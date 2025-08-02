package cmd

import (
	"askai/lib/utils"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Tool for managing configuration file",
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.ReadInConfig(); err != nil {
			cobra.CheckErr(err)
		}
		if err := viper.WriteConfigTo(os.Stdout); err != nil {
			cobra.CheckErr(err)
		}
	},
}

var editConfigCmd = &cobra.Command{
	Use:   "edit",
	Short: "Interactively configure model and API key",
	Run: func(cmd *cobra.Command, args []string) {
		text := utils.Input("Enter provider: ")
		viper.Set("provider", text)

		text = utils.Input("Enter model: ")
		viper.Set("model", text)

		text = utils.Input("Enter API key: ")
		viper.Set("api_key", text)

		text = utils.Input("Enable rate limit (Y/N): ", utils.EmptyNotAllowed())
		if strings.ToLower(text) == "y" {
			text = utils.Input("Enter maximum llm callings per minute: ", utils.PositiveIntegerOnly())
			viper.Set("rate_limit", text)
		}

		if err := viper.WriteConfig(); err != nil {
			cobra.CheckErr(err)
		}
		fmt.Println("Configuration saved")
	},
}

func init() {
	configCmd.AddCommand(showConfigCmd)
	configCmd.AddCommand(editConfigCmd)
}
