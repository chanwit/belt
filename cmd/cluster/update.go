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
	"strings"

	cmdpkg "github.com/chanwit/belt/cmd"
	"github.com/chanwit/belt/util"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update cluster template",
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

		/*
			if len(args) != 0 {
				cluster = args[0]
			}
		*/

		config, err := ioutil.ReadFile(".belt/" + cluster + "/config.yaml")
		if err != nil {
			return err
		}

		data := make(map[string]map[string]string)
		err = yaml.Unmarshal(config, &data)
		if err != nil {
			return err
		}

		// check driver
		driver, err := cmd.Flags().GetString("driver")
		if err != nil {
			return err
		}

		// if driver not specified
		if driver == "" {
			// any first found
			for k := range data {
				driver = k
				break
			}
		}

		def := data[driver]

		defines, err := cmd.Flags().GetStringSlice("define")
		if err != nil {
			return err
		}

		// use args
		// preserve --define key=value for backward compat
		if len(defines) == 0 {
			defines = args
		}

		for _, d := range defines {
			parts := strings.SplitN(d, "=", 2)
			def[parts[0]] = parts[1]
		}

		data[driver] = def
		out, err := yaml.Marshal(data)
		if err != nil {
			return err
		}

		// TODO if file exists
		err = ioutil.WriteFile(".belt/"+cluster+"/config.yaml", out, 0644)
		if err == nil {
			diff := dmp.New()

			diffs := diff.DiffMain(string(config), string(out), true)
			diffs = diff.DiffCleanupSemantic(diffs)
			diffs = diff.DiffCleanupEfficiency(diffs)

			for _, d := range diffs {
				switch d.Type {
				case dmp.DiffInsert:
					fmt.Print("+")
					fmt.Println(d.Text)
				case dmp.DiffDelete:
					fmt.Print("-")
					fmt.Println(d.Text)
				}
			}
		}

		return err
	},
}

func init() {
	cmdpkg.ClusterCmd.AddCommand(updateCmd)

	updateCmd.Flags().String("driver", "", "driver name")
	updateCmd.Flags().StringSlice("define", []string{}, "cluster definition")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
