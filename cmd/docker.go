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

	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

const (
	flagName = "host"
)

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "run docker command remotely",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.ParseFlags(args)

		node := RootCmd.Flag(flagName).Value.String()
		pos := 0
		var err error
		if node != "" {
			for i, a := range args {
				if a == "--"+flagName {
					pos = i + 2
					break
				}
			}
		} else {
			node, err = util.GetActive()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}

		ips := CacheIP()
		nodes := util.Generate(node)

		for _, n := range nodes {
			ip := ips[n]

			cmdArgs := []string{
				"-q",
				"-o",
				"UserKnownHostsFile=/dev/null",
				"-o",
				"StrictHostKeyChecking=no",
				util.DegitalOcean.SSHUser() + "@" + ip,
				"docker",
			}

			cmdArgs = append(cmdArgs, args[pos:]...)
			sshCmd := exec.Command("ssh", cmdArgs...)
			sshCmd.Stdin = os.Stdin
			sshCmd.Stdout = os.Stdout
			sshCmd.Stderr = os.Stderr

			sshCmd.Run()
		}

	},
}

func init() {
	RootCmd.AddCommand(dockerCmd)
	dockerCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		err := cmd.ParseFlags(args)
		if err != nil {
			return
		}
		pos := 0
		for i, a := range args {
			if a == "docker" {
				pos = i + 1
				break
			}
		}

		cmd.Run(cmd, args[pos:])
	})

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dockerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	RootCmd.Flags().String(flagName, "", "use the specific node to control Docker cluster")

	dockerCmd.DisableFlagParsing = true
}
