// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "display logs of a task",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		task := args[0]
		parts := strings.SplitN(task, ".", 2)
		service := parts[0]
		// id := parts[1]

		node := RootCmd.Flag(flagName).Value.String()
		var err error
		if node == "" {
			node, err = util.GetActive()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		ip := GetIP(node)

		cmdArgs := []string{
			"-q",
			"-o",
			"UserKnownHostsFile=/dev/null",
			"-o",
			"StrictHostKeyChecking=no",
			util.DegitalOcean.SSHUser() + "@" + ip,
			"docker", "service", "tasks", service,
		}

		sshCmd := exec.Command("ssh", cmdArgs...)

		bout, err := sshCmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		cid := ""
		targetNode := ""

		lines := strings.Split(strings.TrimSpace(string(bout)), "\n")
		for _, line := range lines {
			parts := strings.Fields(line)
			if parts[1] == task {
				cid = parts[1] + "." + parts[0]
				pos := strings.LastIndex(line, " ")
				targetNode = line[pos+1:]

				break
			}
		}

		if cid != "" && targetNode != "" {
			ip := GetIP(targetNode)
			cmdArgs := []string{
				"-q",
				"-o",
				"UserKnownHostsFile=/dev/null",
				"-o",
				"StrictHostKeyChecking=no",
				util.DegitalOcean.SSHUser() + "@" + ip,
				"docker", "logs", cid,
			}
			sshCmd := exec.Command("ssh", cmdArgs...)
			sshCmd.Stdin = os.Stdin
			sshCmd.Stdout = os.Stdout
			sshCmd.Stderr = os.Stderr
			err := sshCmd.Run()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(logsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
