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
	"os/exec"
	"strings"

	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a set of machines",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Example: `
  To create 10 nodes 512mb each type,
  $ belt create 512mb node[1:10]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please specify parameters")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		size := args[0]

		image := cmd.Flag("image").Value.String()
		if image == "" {
			image = util.DegitalOcean.Image()
		}

		for i := 1; i < len(args); i++ {

			region := cmd.Flag("region").Value.String()

			if region == "" {
				// check prefix: {region}-{node}{num}
				pattern := args[i]
				parts := strings.SplitN(pattern, "-", 2)
				if len(parts) == 2 {
					region = parts[0]
				}
			}

			if region == "" {
				region = util.DegitalOcean.Region()
			}

			boxes := util.Generate(args[i])

			doArgs := []string{
				"-t",
				util.DegitalOcean.AccessToken(),
				"compute",
				"droplet",
				"create",
				"--region",
				region,
				"--ssh-keys",
				util.DegitalOcean.SSHKey(),
				"--image",
				image,
				"--size",
				size,
			}

			cmdExec := exec.Command("doctl", append(doArgs, boxes...)...)
			bout, err := cmdExec.Output()
			if err != nil {
				fmt.Println(err.Error())
				fmt.Print(string(bout))
				return
			}

			ListDroplets()
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
	createCmd.Flags().String("region", "", "override region setting")
	createCmd.Flags().String("image", "", "override image setting")

}
