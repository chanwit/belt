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
	"sync"

	"github.com/chanwit/belt/ssh"
	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

// hostnameCmd represents the hostname command
var hostnameCmd = &cobra.Command{
	Use:   "hostname",
	Short: "correct hostnames",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ips := CacheIP()
		nodes := []string{}
		for i := 0; i < len(args); i++ {
			nodes = append(nodes, util.Generate(args[i])...)
		}

		MAX := 30

		num := len(nodes)
		loop := num / MAX
		rem := num % MAX
		if rem != 0 {
			loop++
		}

		for i := 1; i <= loop; i++ {
			var wg sync.WaitGroup

			for j := 0; j < MAX; j++ {
				offset := (i-1)*MAX + j
				if offset < len(nodes) {
					node := nodes[offset]
					wg.Add(1)
					go func(node string) {
						defer wg.Done()
						ip := ips[node]
						sshcli, err := ssh.NewNativeClient(
							util.DegitalOcean.SSHUser(), ip, util.DegitalOcean.SSHPort(),
							&ssh.Auth{Keys: util.DefaultSSHPrivateKeys()})
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						sshcli.Output("hostname " + node)
						sshcli.Output("service docker restart")

						fmt.Println(node)
					}(node)
				}
			}
			wg.Wait()
		}
	},
}

func init() {
	RootCmd.AddCommand(hostnameCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostnameCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hostnameCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
