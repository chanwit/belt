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
	"errors"
	"fmt"

	"github.com/apcera/libretto/virtualmachine/digitalocean"
	"github.com/chanwit/belt/drivers"
	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

func GetPublicIP(ips *digitalocean.Networks) string {
	for _, ip := range ips.V4 {
		if ip.Type == "public" {
			return ip.IPAddress
		}
	}
	return ""
}

func CacheIP() map[string]string {
	token := util.DegitalOcean.AccessToken()
	resp, err := drivers.GetAllDroplets(token)
	result := make(map[string]string)
	if err != nil {
		return result
	}

	if resp == nil {
		return result
	}

	for _, d := range resp.Droplets {
		result[d.Name] = GetPublicIP(d.Networks)
	}
	return result
}

func GetIP(node string) string {
	token := util.DegitalOcean.AccessToken()
	resp, err := drivers.GetAllDroplets(token)
	if err != nil {
		return ""
	}

	for _, d := range resp.Droplets {
		if d.Name == node {
			return GetPublicIP(d.Networks)
		}
	}
	return ""
}

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "show the IP address for the compute node",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("require a node name as the argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(GetIP(args[0]))
	},
}

func init() {
	RootCmd.AddCommand(ipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
