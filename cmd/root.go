package cmd

import (
	"fmt"
	"os"

	"github.com/j3ssie/jig/core"
	"github.com/spf13/cobra"
)

var options = core.Options{}

var RootCmd = &cobra.Command{
	Use:   "jig",
	Short: fmt.Sprintf("Jig - Jaeles Intput Generator %v by %v\n", core.VERSION, core.AUTHOR),
	Long:  core.Banner(),
}

// Execute main function
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().IntVarP(&options.Concurrency, "concurrency", "c", 5, "Concurrency")
	RootCmd.PersistentFlags().StringVarP(&options.LogFile, "log", "l", "", "Log File")
	RootCmd.PersistentFlags().BoolVar(&options.Debug, "debug", false, "Debug")
	RootCmd.PersistentFlags().BoolVarP(&options.Verbose, "verobse", "v", false, "Verbose output")
	RootCmd.PersistentFlags().BoolVarP(&options.NoOutput, "no-output", "Q", false, "Verbose output")
	//RootCmd.PersistentFlags().BoolVarP(&options.Quite, "quite", "q", false, "Show only essential information")
	RootCmd.PersistentFlags().BoolVarP(&options.Redirect, "redirect", "L", false, "Enable redirect")
	RootCmd.PersistentFlags().BoolVarP(&options.UseChrome, "chrome", "C", false, "Use Chrome headless to send request")

	RootCmd.PersistentFlags().StringSliceVarP(&options.Params, "params", "p", []string{}, "Exclude module name (Multiple -x flags are accepted)")
	RootCmd.PersistentFlags().StringVarP(&options.Input, "input", "i", "", "Input file")
	RootCmd.PersistentFlags().StringVarP(&options.OutputFolder, "Output", "O", "jinp", "Output folder")
	RootCmd.PersistentFlags().StringVarP(&options.Output, "output", "o", "out.txt", "Output File")
	RootCmd.PersistentFlags().BoolVarP(&options.Helper, "list-otype", "P", false, "List all output type")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if !options.Quite {
		if !options.NoBanner {
			core.Banner()
		}
		fmt.Fprintf(os.Stderr, "Jig - Jaeles Intput Generator %v by %v\n", core.VERSION, core.AUTHOR)
	}

	if options.Helper {
		HelpMessage()
		os.Exit(0)
	}

	core.InitLog(&options)
}

// HelpMessage print help message
func HelpMessage() {
	h := "\nUsage:\n jig [mode] [options]\n"
	h += " jig [scan] -h -- Show usage message\n"
	h += "\nSubcommands:\n"
	h += "  jig scan   --  Generate input from Content of list of URLs\n"
	h += "  jig dirb   --  Generate input from Dirbscan result\n"

	h += "\nAvailable Output Type:\n"
	h += `  location   --  Use Location headers as {{.BaseURL}}`

	h += "\n\nExample commands:\n"
	h += `  jig scan -u https://example.com/ -I location`
	fmt.Println(h)
}
