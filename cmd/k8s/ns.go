package k8s

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newNsCmd() *cobra.Command {
	var execCmd = cobra.Command{
		Use:   "set-ns",
		Short: "k8s exec",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a ns argument")
			}
			return nil
		},
		Run: func(c *cobra.Command, args []string) {
			shell := fmt.Sprintf("kubectl config set-context --current --namespace=%s", args[0])
			fmt.Printf("cmd: %s\n", shell)
			cmd := exec.Command("bash", "-c", shell)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			_ = cmd.Run()
		},
		Example: `
		tools k8s set-ns default
		`,
	}

	execCmd.Flags().StringP("shell", "s", "sh", "exec shell")

	return &execCmd
}
