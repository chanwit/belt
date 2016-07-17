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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"

	cmdpkg "github.com/chanwit/belt/cmd"
	"github.com/chanwit/belt/util"
	"github.com/elgs/gojq"
	"github.com/spf13/cobra"
)

// clsCmd represents the cls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list clusters",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		w := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', 0)
		fmt.Fprintf(w, "CLUSTER\tACTIVE\tLEADER\tMASTERS\t#NODES\n")
		infos, err := ioutil.ReadDir(".belt")
		if err != nil {
			return err
		}

		ips := cmdpkg.CacheIP()
		activeCluster, _ := util.GetActiveCluster()
		for _, info := range infos {
			if info.IsDir() {
				if _, err := os.Stat(".belt/" + info.Name() + "/config.yaml"); err == nil {
					cluster := info.Name()
					clusterDisplay := cluster
					if activeCluster == info.Name() {
						clusterDisplay = clusterDisplay + " *"
					}

					host, err := util.GetActiveByCluster(cluster)
					if err != nil {
						fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d / %d\n", clusterDisplay, "-", "-", "-", 0, 0)
						break
					}

					// no active host
					if ips[host] == "" {
						fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d / %d\n", clusterDisplay, "-", "-", "-", 0, 0)
						break
					}

					data, err := cmdpkg.SwarmNodeList(ips[host])
					// DEBUG: fmt.Println(string(data))
					nodes := []interface{}{}
					err = json.Unmarshal(data, &nodes)
					if err != nil {
						return err
					}

					managers := []string{}
					leader := ""
					ready := 0
					for _, node := range nodes {
						p := gojq.NewQuery(node)
						hostname, err := p.QueryToString("Description.Hostname")
						if err != nil {
							panic(err)
						}

						role, err := p.QueryToString("Spec.Role")
						if err != nil {
							panic(err)
						}

						if role == "manager" {
							isLeader, err := p.QueryToBool("ManagerStatus.Leader")
							if err != nil {
								isLeader = false
								// panic(err)
							}

							managers = append(managers, hostname)
							if isLeader {
								leader = hostname
							}
						}

						status, err := p.QueryToString("Status.State")
						if status == "ready" {
							ready++
						}
					}

					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d / %d\n", clusterDisplay, host, leader, strings.Join(managers, ","), ready, len(nodes))
				}
			}
		}

		w.Flush()

		return nil
		// convention over configuration
		// like this
		// .belt/cluster1/config.yml
		// .belt/cluster2/config.yml
		// so what's about active?
		/*
			$ belt cluster ls
			CLUSTER   LEADER  MASTERS            #NODES
			cluster1  node1   node1,node2,node3   30
			cluster2  node2   node1,node2,node3   20
		*/

		// $ belt use cluster1
		// $ belt active cluster1
	},
}

func init() {
	cmdpkg.ClusterCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
