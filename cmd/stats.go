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
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		service := args[0]
		ch := make(chan string)

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

		lines := strings.Split(strings.TrimSpace(string(bout)), "\n")
		for _, line := range lines {
			parts := strings.Fields(line)
			cid := parts[1] + "." + parts[0]
			pos := strings.LastIndex(line, " ")
			targetNode := line[pos+1:]
			ip := GetIP(targetNode)
			go func(ip, cid, targetNode string, ch chan string) {
				for {
					cmdArgs := []string{
						"-q",
						"-o",
						"UserKnownHostsFile=/dev/null",
						"-o",
						"StrictHostKeyChecking=no",
						util.DegitalOcean.SSHUser() + "@" + ip,
						"docker", "stats", cid, "--no-stream",
					}

					sshCmd := exec.Command("ssh", cmdArgs...)
					bout, err := sshCmd.Output()
					if err != nil {
						// ch <- err.Error()
						break
					}

					lines := strings.SplitN(strings.TrimSpace(string(bout)), "\n", 2)
					ch <- lines[1]

					time.Sleep(1 * time.Second)
				}
			}(ip, cid, targetNode, ch)
		}

		w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
		display := make(map[string]string)
		for {
			line := <-ch
			parts := strings.Fields(line)
			display[parts[0]] = strings.TrimSpace(line[len(parts[0]):])

			keys := []string{}
			for k := range display {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			// clear screen
			fmt.Print("\033[2J")
			fmt.Print("\033[H")

			fmt.Fprint(w, "CONTAINER\tCPU %\tMEM USAGE / LIMIT\tMEM %\tNET I/O\tBLOCK I/O\tPIDS\n")
			for _, k := range keys {
				v := display[k]
				f := strings.Fields(v)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s",
					k,
					f[0],
					strings.Join(f[1:5], " "),
					f[6],
					strings.Join(f[7:11], " "),
					strings.Join(f[12:16], " "),
					f[17],
				)
				fmt.Fprintln(w)
			}
			w.Flush()
		}

	},
}

func init() {
	RootCmd.AddCommand(statsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
