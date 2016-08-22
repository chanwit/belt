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

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := util.GetActiveCluster()
		if err != nil {
			return err
		}

		activeNode, err := util.GetActive()
		if err != nil {
			return err
		}

		if val, _ := cmd.Flags().GetBool("unset"); val {
			fmt.Printf("unset MACHINE_STORAGE_PATH\n")
			machineCmd := exec.Command("docker-machine",
				"env",
				"-u",
			)
			machineCmd.Stdin = os.Stdin
			machineCmd.Stdout = os.Stdout
			machineCmd.Stderr = os.Stderr
			return machineCmd.Run()
		} else {
			beltMachinePath := ".belt/" + cluster + "/machine"
			fmt.Printf("export MACHINE_STORAGE_PATH=\"%s\"\n", beltMachinePath)
			machineCmd := exec.Command("docker-machine",
				"-s",
				beltMachinePath,
				"env",
				activeNode,
			)
			machineCmd.Stdin = os.Stdin
			machineCmd.Stdout = os.Stdout
			machineCmd.Stderr = os.Stderr
			return machineCmd.Run()
		}
	},
}

func init() {
	RootCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	envCmd.Flags().BoolP("unset", "u", false, "Unset environment variables")
}
