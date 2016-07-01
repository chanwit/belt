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
	"net"
	"os"
	"os/exec"
	"sync"

	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		master := args[0]
		if net.ParseIP(master) == nil {
			master = GetIP(master)
		}

		_ /*pwd*/, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		nodes := util.Generate(args[1])

		num := len(nodes)
		loop := num / 10
		rem := num % 10
		if rem != 0 {
			loop++
		}

		for i := 1; i <= loop; i++ {
			var wg sync.WaitGroup

			for j := 0; j < 10; j++ {
				offset := (i-1)*10 + j
				if offset < len(nodes) {
					node := nodes[offset]
					wg.Add(1)
					go func(node string) {
						defer wg.Done()
						ip := GetIP(node)
						sshCmd := exec.Command("ssh",
							"-q",
							"-o",
							"UserKnownHostsFile=/dev/null",
							"-o",
							"StrictHostKeyChecking=no",
							util.DegitalOcean.SSHUser()+"@"+ip,
							"docker", "swarm", "join", master+":2377",
						)
						bout, err := sshCmd.CombinedOutput()
						if err != nil {
							fmt.Println(err.Error())
						}
						fmt.Print(node + ": " + string(bout))
					}(node)
				}
			}

			wg.Wait()

		}

	},
}

func init() {
	swarmCmd.AddCommand(joinCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// joinCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// joinCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
