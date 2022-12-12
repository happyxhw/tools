package k8s

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	ns          string
	resourceMap = map[string]string{
		"deploy": "deployment",
		"svc":    "svc",
		"pod":    "pod",
	}
	actionMap = map[string]string{
		"get":  "get",
		"desc": "describe",
	}
)

func NewCmd() *cobra.Command {
	var k8sCmd = cobra.Command{
		Use:   "k8s",
		Short: "k8s tools",
		Long:  `k8s tools`,
		Run: func(c_md *cobra.Command, _ []string) {
			if len(os.Args) <= 2 {
				return
			}
			if len(os.Args) >= 5 {
				action := actionMap[os.Args[2]]
				if action != "" {
					r := resourceMap[os.Args[3]]
					if r == "" {
						return
					}
					pod := searchResource(os.Args[4], ns, r)
					if pod == "" {
						return
					}
					var args []string
					if ns != "" {
						args = append(args, []string{"-n", ns}...)
					}
					args = append(args, []string{action, r, pod}...)
					if len(os.Args) >= 6 {
						for i := 5; i < len(os.Args); i++ {
							if os.Args[i] == "--" {
								continue
							}
							args = append(args, os.Args[i])
						}
					}

					cmd := exec.Command("kubectl", args...)
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					fmt.Println("cmd: ", cmd.String())
					_ = cmd.Run()
					return
				}
			}

			args := os.Args[2:len(os.Args)]
			cmd := exec.Command("kubectl", args...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			fmt.Println("cmd: ", cmd.String())
			_ = cmd.Run()
		},
	}
	k8sCmd.Flags().StringVarP(&ns, "ns", "n", "", "k8s namespace")

	k8sCmd.AddCommand(newNsCmd())
	k8sCmd.AddCommand(newExecCmd())
	k8sCmd.AddCommand(newLogsCmd())

	return &k8sCmd
}

func searchResource(pod, ns, resource string) string {
	shell := fmt.Sprintf(
		"%s get %s | grep -e '%s' | awk '{print $1}'", kubectl(ns), resource, pod,
	)

	cmd := exec.Command("bash", "-c", shell)
	fmt.Printf("cmd: %s\n", shell)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("search pod: %s", err)
		os.Exit(-1)
	}
	if len(out) == 0 {
		return ""
	}
	outStr := strings.TrimSuffix(string(out), "\n")
	pods := strings.Split(outStr, "\n")
	for i, p := range pods {
		fmt.Printf("pod-%d: %s\n", i+1, p)
	}
	index := 1
	if len(pods) > 1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Please Select Pod Index, 1~%d: ", len(pods))
		indexStr, err := reader.ReadString('\n')
		if err != nil {
			indexStr = ""
		}
		indexStr = strings.TrimSuffix(indexStr, "\n")
		index, err = strconv.Atoi(indexStr)
		if err != nil || index > len(pods) || index <= 0 {
			index = 1
		}
	}

	return pods[index-1]
}

func kubectl(ns string) string {
	c := "kubectl"
	if ns != "" {
		c = fmt.Sprintf("%s -n %s", c, ns)
	}
	return c
}
