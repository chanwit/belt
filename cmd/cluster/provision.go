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
	"strings"
	"sync"
	"time"

	cmdpkg "github.com/chanwit/belt/cmd"
	"github.com/chanwit/belt/drivers"
	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

func stringSliceConv(sval string) ([]string, error) {
	sval = strings.TrimPrefix(sval, "[")
	sval = strings.TrimSuffix(sval, "]")
	// An empty string would cause a slice with one (empty) string
	if len(sval) == 0 {
		return []string{}, nil
	}
	v := strings.Split(sval, ",")
	return v, nil
}

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
		masters, err := stringSliceConv(cmd.Flag("manager").Value.String())

		masterSize, err := cmd.Flags().GetString("manager-size")
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

		secret, err := cmd.Flags().GetString("secret")
		if err != nil {
			return err
		}

		if secret == "" {
			return fmt.Errorf("secret must be specified")
		}

		genMasters := []string{}
		for _, m := range masters {
			genMasters = append(genMasters, util.Generate(m)...)
		}

		// fmt.Printf("%q\n",genMasters)

		// 1. create masters
		masterConfig := drivers.Config{
			Names:   genMasters,
			Region:  util.DegitalOcean.Region(),
			Image:   util.DegitalOcean.Image(),
			Size:    masterSize,
			SSHKeys: []string{util.DegitalOcean.SSHKey()},
		}

		masterDroplets, err := drivers.Provision(util.DegitalOcean.AccessToken(), masterConfig)
		if err != nil {
			return err
		}

		// 4. create worker
		workers := []string{}
		for _, arg := range args {
			workers = append(workers, util.Generate(arg)...)
		}

		workerConfig := drivers.Config{
			Names:   workers,
			Region:  util.DegitalOcean.Region(),
			Image:   util.DegitalOcean.Image(),
			Size:    workerSize,
			SSHKeys: []string{util.DegitalOcean.SSHKey()},
		}

		workerDroplets, err := drivers.Provision(util.DegitalOcean.AccessToken(), workerConfig)
		if err != nil {
			return err
		}

		// print out masters first
		cmdpkg.ListDroplets(masterDroplets.Droplets)

		fmt.Print("\nwaiting for all masters to be active ...")
		// 2. wait for masters to be active
		status := make(map[string]bool)
		for {
			res, _ := drivers.GetAllDroplets(util.DegitalOcean.AccessToken())
			for _, d := range res.Droplets {
				//fmt.Printf("%s = %s\n", d.Name, d.Status)
				status[d.Name] = (d.Status == "active")
			}

			allActive := true
			for _, m := range genMasters {
				if status[m] == false {
					allActive = false
					break
				}
			}

			if allActive {
				break
			} else {
				fmt.Print(".")
				time.Sleep(3 * time.Second)
			}

		}

		fmt.Println()

		// 3. init and join
		// if len(master) == 1 // do init
		// if len(master) >= 2 // do init, do join --manager

		util.SetActive(genMasters[0])

		fmt.Println("initialising a cluster ...")
		// swarm init
		ips := cmdpkg.CacheIP()
		primeIP := ips[genMasters[0]]
		sout, err := cmdpkg.SwarmInit(primeIP, secret)
		if err != nil {
			return err
		}
		// check CA hash
		fmt.Println(sout)

		fmt.Println(genMasters[0] + ": init ...")

		// todo handle error
		cmdpkg.SwarmUpdate(primeIP, "manager")
		defer func() {
			_, err := cmdpkg.SwarmUpdate(primeIP, "none")
			if err == nil {
				fmt.Println("acceptance policy updated to none")
			}
		}()

		fmt.Println(genMasters[0] + ": policy updated")

		var wg sync.WaitGroup
		for _, m := range genMasters[1:] {
			wg.Add(1)
			go func(m string) {
				defer wg.Done()
				ip := ips[m]
				_, err := cmdpkg.SwarmJoinAsMaster(ip, primeIP, secret)
				if err != nil {
					panic(err)
				}
				fmt.Println(m + ": joined as manager")
			}(m)
		}
		wg.Wait()

		// print out workers
		cmdpkg.ListDroplets(workerDroplets.Droplets)

		// 5. wait for worker to be active
		for {
			res, _ := drivers.GetAllDroplets(util.DegitalOcean.AccessToken())
			for _, d := range res.Droplets {
				status[d.Name] = (d.Status == "active")
			}

			allActive := true
			count := 0
			for _, w := range workers {
				if status[w] == false {
					allActive = false
				} else {
					count++
				}
			}

			if allActive {
				break
			} else {
				fmt.Printf("\r%d / %d worker nodes become active ...", count, len(workers))
				// should reduce from 10 -> 5 -> 3
				// when most nodes done
				time.Sleep(5 * time.Second)
			}

		}

		fmt.Printf("\r%d / %d worker nodes become active ...\n", len(workers), len(workers))

		// update policy to accept worker
		// todo handle error
		cmdpkg.SwarmUpdate(primeIP, "worker")
		fmt.Println("acceptance policy updated to worker")

		// 5. join as worker
		for _, w := range workers {
			wg.Add(1)
			go func(w string) {
				defer wg.Done()
				ip := ips[w]
				_, err := cmdpkg.SwarmJoinAsWorker(ip, primeIP, secret)
				if err != nil {
					panic(err)
				}
				fmt.Println(w + ": joined as worker")
			}(w)
		}
		wg.Wait()

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
	provisionCmd.Flags().String("manager-size", "", "manager size")
	provisionCmd.Flags().String("worker-size", "", "worker size")

	provisionCmd.Flags().StringSlice("manager", []string{}, "managers")
	// provisionCmd.Flags().StringSlice("worker", []string{}, "workers")
	provisionCmd.Flags().String("secret", "", "secret for forming cluster")
}
