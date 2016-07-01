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
	"os"
	"os/exec"
	// "strings"

	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "delete compute nodes",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		boxes := util.Generate(args[0])
		boxesMap := make(map[string]bool)
		for _, box := range boxes {
			boxesMap[box] = true
		}

		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		lsArgs := []string{
			"-c",
			pwd + "/.doctlcfg",
			"-o",
			"json",
			"compute",
			"droplet",
			"ls",
		}

		lsExec := exec.Command("doctl", lsArgs...)
		bout, err := lsExec.Output()
		nodes := []status{}
		err = json.Unmarshal(bout, &nodes)

		boxToRm := []string{}
		for _, node := range nodes {
			if node.Status == "active" && boxesMap[node.Name] == true {
				fmt.Println(node.Name)
				boxToRm = append(boxToRm, node.Name)
			}
		}

		doArgs := []string{
			"-t",
			util.DegitalOcean.AccessToken(),
			"compute",
			"droplet",
			"rm",
		}

		cmdExec := exec.Command("doctl", append(doArgs, boxToRm...)...)
		rmOut, err := cmdExec.CombinedOutput()
		if err != nil {
			fmt.Println(err.Error())
		}

		// print output from err, if exists
		fmt.Print(string(rmOut))
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
