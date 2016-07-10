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
	"strings"

	"github.com/chanwit/belt/ssh"
	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove nodes from the swarm",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		nodes := []string{}
		for i := 0; i < len(args); i++ {
			nodes = append(nodes, util.Generate(args[i])...)
		}

		host, err := util.GetActive()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		ip := GetIP(host)
		sshcli, err := ssh.NewNativeClient(
			util.DegitalOcean.SSHUser(), ip, util.DegitalOcean.SSHPort(),
			&ssh.Auth{Keys: util.DefaultSSHPrivateKeys()})

		if err != nil {
			fmt.Println(err.Error())
		}

		sout, err := sshcli.Output("docker node remove " + strings.Join(nodes, " "))
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Print(sout)
	},
}

func init() {
	swarmCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
