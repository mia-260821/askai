package cmd

import (
	"askai/lib/utils"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path"
)

const ConfigFile = ".askai.yaml"

var (
	cfgFile string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config-file", "", fmt.Sprintf("config file (default is $HOME/%s)", ConfigFile))
	// add subcommands
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(chatCmd)
}

func initConfig() error {
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		cfgFile = path.Join(home, ConfigFile)
	}
	if err := utils.FileCreatesIfNotExists(cfgFile); err != nil {
		return err
	}
	viper.SetConfigFile(cfgFile)
	return viper.ReadInConfig()
}

var rootCmd = &cobra.Command{
	Use:   "askai",
	Short: "askai is a sample CLI tool written in Go",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := initConfig(); err != nil {
			cobra.CheckErr(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from AskAI!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
