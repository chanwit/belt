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
	"bytes"
	"fmt"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/chanwit/belt/ssh"
	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

var configTmpl = `
[global_tags]

[agent]
  interval = "10s"
  round_interval = true
  metric_batch_size = 1000
  metric_buffer_limit = 10000
  collection_jitter = "0s"
  flush_interval = "10s"
  flush_jitter = "0s"
  precision = ""
  debug = false
  quiet = false
  hostname = ""
  omit_hostname = false

[[outputs.influxdb]]
  urls = ["http://{{.host}}:8086"] # required
  database = "telegraf" # required
  retention_policy = ""
  write_consistency = "any"
  timeout = "5s"
  username = "{{.user}}"
  password = "{{.pass}}"

[[inputs.cpu]]
  percpu = true
  totalcpu = true
  fielddrop = ["time_*"]

[[inputs.mem]]

[[inputs.conntrack]]
  files = ["ip_conntrack_count","ip_conntrack_max",
            "nf_conntrack_count","nf_conntrack_max"]
  dirs = ["/proc/sys/net/ipv4/netfilter","/proc/sys/net/netfilter"]

[[inputs.docker]]
  endpoint = "unix:///var/run/docker.sock"
  timeout = "5s"
`

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		dbhost := cmd.Flag("dbhost").Value.String()
		if dbhost == "" {
			fmt.Println("dbhost is required")
			return
		}

		ips := CacheIP()
		nodes := []string{}
		for _, arg := range args {
			nodes = append(nodes, util.Generate(arg)...)
		}

		// generate config from template
		tmpl, err := template.New("config").Parse(configTmpl)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data := map[string]string{
			"host": ips[dbhost],
			"user": "root",
			"pass": "root",
		}
		sbuffer := bytes.NewBufferString("")
		err = tmpl.Execute(sbuffer, data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var wg sync.WaitGroup
		for _, n := range nodes {
			wg.Add(1)
			go func(node string) {
				defer wg.Done()

				ip := ips[node]

				sshcli, err := ssh.NewNativeClient(
					util.DegitalOcean.SSHUser(), ip, util.DegitalOcean.SSHPort(),
					&ssh.Auth{Keys: util.DefaultSSHPrivateKeys()})
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				// pre-check
				var telegraf bool
				sout, err := sshcli.Output("/usr/bin/telegraf version")
				if err != nil {
					telegraf = false
				}

				if strings.TrimSpace(sout) == "Telegraf - version 1.0.0-beta2" {
					telegraf = true
				} else {
					// need upgrade
					sshcli.Output("apt-get remove -y telegraf")
					sshcli.Output("rm /etc/telegraf/telegraf.conf")

					telegraf = false
				}

				// setup agent
				if !telegraf {
					fmt.Println(node + ": installing agent ...")
					sshcli.Output("wget -O /tmp/telegraf.deb https://dl.influxdata.com/telegraf/releases/telegraf_1.0.0-beta2_amd64.deb")
					// TODO
					sshcli.Output("dpkg -i /tmp/telegraf.deb")
					sshcli.Output(fmt.Sprintf("echo '%s' | tee /etc/telegraf/telegraf.conf", sbuffer.String()))

					// post check
					for {
						_, err := sshcli.Output("service telegraf status")
						if err != nil {
							sshcli.Output("service telegraf start")
						} else {
							fmt.Println(node + ": done")
							break
						}
						time.Sleep(500 * time.Millisecond)
					}
				} else {
					fmt.Println(node + ": telegraf already installed")
				}
			}(n)
		}
		wg.Wait()

		// setup: monitor node[1:50] --dashboard mon
		// setup: monitor node[2:30]
		// if --dashboard
		// 1. setup influxsrv
		// 2. setup grafana
		// loop over machines
		// pre-check
		//   1. check version
		//   if equals
		// 3. wget telegraf
		// 4. gen config
		// 5. upload config to node /etc/...
		// 6. dpkg -i telegraf
		// post-check
	},
}

func init() {
	RootCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	monitorCmd.Flags().String("dbhost", "", "hostname for monitoring database (InfluxDB)")
	// TODO dbuser, dbpass
}
