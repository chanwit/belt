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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		size := args[0]
		boxes := util.Generate(args[1])
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		doArgs := []string{
			"-c",
			pwd + "/.doctlcfg",
			"compute",
			"droplet",
			"create",
			"--region",
			"sgp1",
			"--ssh-keys",
			"816630",
			"--image",
			"18153887",
			"--size",
			size,
		}

		cmdExec := exec.Command("doctl", append(doArgs, boxes...)...)
		bout, err := cmdExec.Output()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		lines := strings.Split(string(bout), "\n")
		for i, line := range lines {
			if i == 0 || i%2 == 1 {
				fmt.Println(line)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
