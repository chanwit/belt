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

	"github.com/spf13/cobra"
	cmdpkg "github.com/chanwit/belt/cmd"
	"github.com/chanwit/belt/util"
	"github.com/chanwit/belt/drivers"
)

// provisionCmd represents the provision command
var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "provision cluster and form a swarm",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		masters, err := cmd.Flags().GetStringSlice("master")

		masterSize, err := cmd.Flags().GetString("master-size")
		workerSize, err := cmd.Flags().GetString("worker-size")

		if masterSize == "" {
			masterSize = util.DegitalOcean.Size()
		}
		if workerSize == "" {
			workerSize = util.DegitalOcean.Size()
		}

		if workerSize == "" || masterSize == "" {
			return fmt.Errorf("size must be specified.")
		}

		masterConfig := drivers.Config{
    		Names:  masters,
    		Region: util.DegitalOcean.Region(),
    		Image:  util.DegitalOcean.Image(),
    		Size:   masterSize,
		}

		masterDroplets, err := drivers.Provision(util.DegitalOcean.AccessToken(), masterConfig);
		if err != nil {
			return err
		}

		cmdpkg.ListDroplets(masterDroplets.Droplets)

		status := make(map[string]bool)
		for {
			res, _ := drivers.GetAllDroplets(util.DegitalOcean.AccessToken())
			for _, d := range res.Droplets {
				status[d.Name] = (d.Status == "active")
			}

			/*
			allActive = false
			for _, m := range masters {

			}
			*/
		}

		// create masters
		// create worker
		// wait for masters to be active
		// if len(master) == 1 // do init
		// if len(master) >= 2 // do init, do join --manager
		return nil
	},
}

func init() {
	cmdpkg.ClusterCmd.AddCommand(provisionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// provisionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	provisionCmd.Flags().String("master-size", "", "master size")
	provisionCmd.Flags().String("worker-size", "", "worker size")

	provisionCmd.Flags().StringSlice("master", []string{}, "masters")
	provisionCmd.Flags().StringSlice("worker", []string{}, "workers")
}
