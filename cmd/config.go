package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Print("Enter provider: ")
		if scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())
			viper.Set("provider", text)
		}
		fmt.Print("Enter model: ")
		if scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())
			viper.Set("model", text)
		}
		fmt.Print("Enter API key: ")
		if scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())
			viper.Set("api_key", text)
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
