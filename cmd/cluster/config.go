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

package cluster

import (
	"fmt"
	"io/ioutil"

	cmdpkg "github.com/chanwit/belt/cmd"
	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "show the current configuration",
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
		b, err := ioutil.ReadFile(".belt/" + cluster + "/config.yaml")
		if err != nil {
			return err
		}
		showAsCmd, err := cmd.Flags().GetBool("cmd")
		if showAsCmd == false {
			fmt.Print(string(b))
		} else {
			data := make(map[string]map[string]string)
			yaml.Unmarshal(b, &data)
			for k, def := range data {
				fmt.Printf("--driver %s ", k)
				for kk, vv := range def {
					fmt.Printf("--define %s=%s ", kk, vv)
				}
				// only once
				break
			}
		}
		return nil
	},
}

func init() {
	cmdpkg.ClusterCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	configCmd.Flags().BoolP("cmd", "c", false, "print as command line")

}
