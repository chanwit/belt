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
	"os"
	"strings"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
	cmdpkg "github.com/chanwit/belt/cmd"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new cluster",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// $ belt cluster new --name=cluster \
		//   --driver digitalocean \
		//     --define region=x \
		//     --define image=y \
		//     --define access-token=abc
		os.MkdirAll(".belt/" + args[0], 0755)

		def := make(map[string]string)
		defines, err := cmd.Flags().GetStringSlice("define")
		if err != nil {
			return err
		}

		for _, d := range defines {
			parts := strings.SplitN(d, "=", 2)
			def[parts[0]] = parts[1]
		}

		driver, err := cmd.Flags().GetString("driver")
		if err != nil {
			return err
		}

		data := make(map[string]interface{})
		data[driver] = def
		out, err := yaml.Marshal(data)
		if err != nil {
			return err
		}

		// TODO if file exists
		return ioutil.WriteFile(".belt/" + args[0] + "/config.yaml", out, 0644)
	},
}

func init() {
	cmdpkg.ClusterCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	newCmd.Flags().String("name", "", "name of a new cluster")
	newCmd.Flags().String("driver", "none", "driver name")
	newCmd.Flags().StringSlice("define", []string{}, "cluster definition")
}
