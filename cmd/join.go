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
	"strings"
	"sync"

	"github.com/chanwit/belt/ssh"
	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

func SwarmUpdate(ip string, policy string) (string, error) {
	_, err := GetSSHClient(ip)
	if err != nil {
		return "", err
	}

	// return sshcli.Output("docker swarm update --auto-accept " + policy)
	return "", nil
}

func SwarmJoinAsMaster(ip string, prime string, token string) (string, error) {
	sshcli, err := GetSSHClient(ip)
	if err != nil {
		return "", err
	}

	return sshcli.Output(fmt.Sprintf("docker swarm join --listen-addr %s:2377 --advertise-addr %s:2377 --token %s %s", ip, ip, token, prime))
}

func SwarmJoinAsWorker(ip string, prime string, token string) (string, error) {
	sshcli, err := GetSSHClient(ip)
	if err != nil {
		return "", err
	}

	return sshcli.Output(fmt.Sprintf("docker swarm join --listen-addr %s:2377 --advertise-addr %s:2377 --token %s %s", ip, ip, token, prime))
}

func SwarmToken(ip string, kind string) (string, error) {
	sshcli, err := GetSSHClient(ip)
	if err != nil {
		return "", err
	}

	return sshcli.Output(fmt.Sprintf("docker swarm join-token -q %s", kind))
}

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		prime, err := util.GetActive()
		if err != nil {
			return err
		}

		isManager, err := cmd.Flags().GetBool("manager")
		if err != nil {
			return err
		}

		ips := CacheIP()
		if net.ParseIP(prime) == nil {
			prime = ips[prime]
		}

		token := ""
		if isManager {
			token, err = SwarmToken(prime, "manager")
			if err != nil {
				return err
			}
		} else {
			token, err = SwarmToken(prime, "worker")
			if err != nil {
				return err
			}
		}

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

						cmd := "docker swarm join --listen-addr %s:2377 --advertise-addr %s:2377 --token %s %s:2377"
						cmd = fmt.Sprintf(cmd, ip, ip, strings.TrimSpace(token), prime)
						sout, err := sshcli.Output(cmd)
						if err != nil {
							fmt.Println(err.Error())
						}

						fmt.Print(node + ": " + sout)
					}(node)
				}
			}

			wg.Wait()

		}

		return nil

		// set
	},
}

func init() {
	RootCmd.AddCommand(joinCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// joinCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// joinCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	joinCmd.Flags().Bool("manager", false, "join as manager")
	// joinCmd.Flags().BoolP("enable-remote", "m", false, "allow remote connection to Engine")
	// joinCmd.Flags().StringP("secret", "s", "", "secret for cluster")
}
