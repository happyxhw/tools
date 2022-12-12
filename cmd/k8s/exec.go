package k8s

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newExecCmd() *cobra.Command {
	var execCmd = cobra.Command{
		Use:   "exec",
		Short: "k8s exec",
		Long:  ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a pod argument")
			}
			return nil
		},
		Run: func(c *cobra.Command, args []string) {
			pod := searchPod(args[0], ns)
			if pod == "" {
				return
			}
			shell := fmt.Sprintf("%s exec -it %s", kubectl(ns), pod)
			if sh := c.Flag("shell").Value.String(); sh != "" {
				shell += " -- " + sh
			}
			fmt.Printf("cmd: %s\n", shell)
			cmd := exec.Command("bash", "-c", shell)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			fmt.Println()
			_ = cmd.Run()
		},
		Example: `
		tools k8s exec pod_like -n ns -s bash
		`,
	}

	execCmd.Flags().StringP("shell", "s", "sh", "exec shell")

	return &execCmd
}
