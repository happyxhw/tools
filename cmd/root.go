package cmd

/*
doc: https://zetcode.com/golang/exec-command/
*/

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/happyxhw/tools/cmd/k8s"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short: "tools",
	Long:  "tools",
	Run: func(cmd *cobra.Command, args []string) {
		debug()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(k8s.NewCmd())
}

func debug() {
	fmt.Println("debug")
}
