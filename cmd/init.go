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

	"github.com/chanwit/belt/ssh"
	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

var clients map[string]ssh.Client

func init() {
	clients = make(map[string]ssh.Client)
}

func ClearSSHClient(ip string) {
	delete(clients, ip)
}

func GetSSHClient(ip string) (ssh.Client, error) {
	cli, exist := clients[ip]
	if exist {
		return cli, nil
	}
	sshcli, err := ssh.NewNativeClient(
		util.DegitalOcean.SSHUser(), ip, util.DegitalOcean.SSHPort(),
		&ssh.Auth{Keys: util.DefaultSSHPrivateKeys()})

	if err != nil {
		return nil, err
	}

	clients[ip] = sshcli
	return sshcli, nil
}

// belt docker node update --availability drain mg0 mg1 mg2^C

func DrainNodes(ip string, nodes []string) (string, error) {
	sshcli, err := GetSSHClient(ip)
	if err != nil {
		return "", err
	}

	// docker node update --availability drain mg0 mg1 mg2
	result := []string{}
	for _, node := range nodes {
		sout, err := sshcli.Output("docker node update --availability drain " + node)
		if err != nil {
			fmt.Print(sout)
		}
		result = append(result, strings.TrimSpace(sout))
	}

	return strings.Join(result, "\n"), err
}

func SwarmInit(ip string) (string, error) {
	sshcli, err := GetSSHClient(ip)
	if err != nil {
		return "", err
	}

	// use accept none
	sout, err := sshcli.Output("docker swarm init --listen-addr " + ip + ":2377 --advertise-addr " + ip + ":2377")
	if err != nil {
		fmt.Print(sout)
	}
	return sout, err
}

func SwarmNodeList(ip string) ([]byte, error) {
	sshcli, err := GetSSHClient(ip)
	if err != nil {
		return nil, err
	}

	// use accept none
	sout, err := sshcli.Output("curl -s --unix-socket /var/run/docker.sock http:/v1.24/nodes")
	if err != nil {
		fmt.Print(sout)
	}
	return []byte(sout), err
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a swarm of Docker Engines",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please specify machines")
		}

		node := args[0]
		ip := GetIP(node)

		cluster, err := util.GetActiveCluster()
		if err != nil {
			return err
		}

		beltMachinePath := ".belt/" + cluster + "/machine"
		if err != nil {
			return err
		}

		if val, _ := cmd.Flags().GetBool("tls"); val {

			machineCmd := exec.Command("docker-machine",
				"-s",
				beltMachinePath,
				"create",
				"-d", "generic",
				"--generic-ip-address="+ip,
				node,
			)
			machineCmd.Stdin = os.Stdin
			// machineCmd.Stdout = os.Stdout
			machineCmd.Stderr = os.Stderr
			fmt.Println("Enabling TLS with docker-machine ...")
			err := machineCmd.Run()
			if err != nil {
				fmt.Println("Error executing docker-machine: " + err.Error())
				return err
			}
		}

		sout, err := SwarmInit(ip)
		if err != nil {
			fmt.Print(sout)
			return err
		}

		util.SetActive(node)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().Bool("tls", true, "setup TLS for the Engine")
}
