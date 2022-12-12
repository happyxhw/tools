package k8s

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	tail   int
	follow bool
)

func newLogsCmd() *cobra.Command {
	var execCmd = cobra.Command{
		Use:   "logs",
		Short: "k8s logs",
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
			shell := fmt.Sprintf("%s logs --tail=%d", kubectl(ns), tail)
			if follow {
				shell += " -f"
			}
			shell += " " + pod
			fmt.Printf("cmd: %s\n", shell)
			cmd := exec.Command("bash", "-c", shell)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			fmt.Println()
			_ = cmd.Run()
		},
		Example: `
		tools k8s logs pod_like -f --tail=100
		`,
	}

	execCmd.Flags().StringP("shell", "s", "sh", "exec shell")
	execCmd.Flags().IntVarP(&tail, "tail", "t", 20, "logs tail count")
	execCmd.Flags().BoolVarP(&follow, "follow", "f", false, "logs follow")

	return &execCmd
}
