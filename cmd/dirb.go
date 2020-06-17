package cmd

import (
	"fmt"

	"github.com/j3ssie/jig/core"
	"github.com/spf13/cobra"
)

func init() {
	var configCmd = &cobra.Command{
		Use:   "dirb",
		Short: "Generate input from DirbScan module",
		Long:  core.Banner(),
		RunE:  runDirbb,
	}

	configCmd.Flags().StringP("action", "a", "select", "Action")
	configCmd.Flags().StringP("pluginsRepo", "p", "git@gitlab.com:j3ssie/osmedeus-plugins.git", "Osmedeus Plugins repository")
	// for cred action
	configCmd.Flags().String("user", "", "Username")
	configCmd.Flags().String("pass", "", "Password")
	configCmd.Flags().StringP("workspace", "w", "", "Name of workspace")

	//configCmd.SetHelpFunc(configHelp)
	RootCmd.AddCommand(configCmd)
}

func runDirbb(cmd *cobra.Command, _ []string) error {
	// action, _ := cmd.Flags().GetString("action")

	fmt.Println("Dirb command")

	return nil
}
