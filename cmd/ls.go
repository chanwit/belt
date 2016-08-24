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
	"text/tabwriter"

	"github.com/apcera/libretto/virtualmachine/digitalocean"
	"github.com/chanwit/belt/drivers"
	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

func ListDroplets() {
	token := util.DegitalOcean.AccessToken()
	resp, err := drivers.GetAllDroplets(token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if resp == nil {
		PrintDroplets([]*digitalocean.Droplet{})
	} else {
		PrintDroplets(resp.Droplets)
	}
}

func PrintDroplets(droplets []*digitalocean.Droplet) {
	w := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', 0)
	fmt.Fprintf(w, "NAME\tIPv4\tMEMORY\tREGION\tIMAGE\tSTATUS\n")
	if droplets != nil {
		for _, d := range droplets {
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\n",
				d.Name,
				GetPublicIP(d.Networks),
				d.Size.Memory,
				d.Region.Slug,
				d.Image.Distribution+" "+d.Image.Name,
				d.Status)
		}
	}
	w.Flush()
}

// llsCmd represents the lls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list machines",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ListDroplets()
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// llsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// llsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
