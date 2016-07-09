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
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

type status struct {
	Name   string
	Status string
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "show status of compute nodes",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		wait := cmd.Flag("wait").Value.String()
		what := ""
		num := 0
		// name := ""

		if wait != "" {
			// normal format
			// active=5
			parts := strings.SplitN(wait, "=", 2)
			what = parts[0]

			if what != "new" && what != "active" {
				// TODO
				// status to check
				// this format: mon=active
				// what = parts[1]
				// name = parts[0]
			} else {
				var err error
				num, err = strconv.Atoi(parts[1])
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}

		doArgs := []string{
			"-t",
			util.DegitalOcean.AccessToken(),
			"-o",
			"json",
			"compute",
			"droplet",
			"ls",
		}

		done := false
		loop := false

		for {

			cmdExec := exec.Command("doctl", doArgs...)
			bout, err := cmdExec.Output()
			nodes := []status{}
			err = json.Unmarshal(bout, &nodes)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// group by status
			grpStatus := make(map[string][]status)
			grpStatus["new"] = []status{}
			grpStatus["active"] = []status{}
			for _, node := range nodes {
				grpStatus[node.Status] = append(grpStatus[node.Status], node)
			}

			if loop {
				fmt.Println()
			}
			fmt.Printf("%s\t\t%s\t%s\n", "STATUS", "#NODES", "NAMES")

			MAX_NAMES_LEN := 56
			for k, v := range grpStatus {
				if len(v) != 0 {
					names := v[0].Name
					for i, vv := range v {
						if i != 0 {
							names = names + ", " + vv.Name
						}
						if len(names) > MAX_NAMES_LEN {
							names = names[0:MAX_NAMES_LEN] + " ..."
							break
						}
					}
					fmt.Printf("%s\t\t%5d\t%s\n", k, len(v), names)
				}

				if k == what && len(v) == num {
					done = true
				}
			}

			if wait == "" {
				done = true
			}

			if done {
				return
			}

			time.Sleep(10 * time.Second)
			loop = true
		}

	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringP("wait", "w", "", "wait until the criteria match")
}
