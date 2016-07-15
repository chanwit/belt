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
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/chanwit/belt/ssh"
	"github.com/chanwit/belt/util"
	"github.com/spf13/cobra"
)

const dataSource = `
{
	"name":"telegraf",
	"type":"influxdb",
	"access":"proxy",
	"url":"http://{{.host}}:8086",
	"user":"root",
	"password":"root",
	"database":"telegraf",
	"basicAuth":true,
	"basicAuthUser":"admin",
	"basicAuthPassword":"admin",
	"withCredentials":false,
	"isDefault":true
}
`

func CheckDashboard(node string) bool {
	client := &http.Client{}
	req, err := buildRequest("GET", "http://"+GetIP(node)+"/api/dashboards/db/belt", nil)
	if err != nil {
		return false
	}
	rsp, err := client.Do(req)
	if err != nil {
		return false
	}

	// OK
	if rsp.StatusCode >= 200 && rsp.StatusCode <= 299 {
		return true
	}

	return false
}

func InstallDashboard(node string) error {
	ip := GetIP(node)
	sshcli, err := ssh.NewNativeClient(
		util.DegitalOcean.SSHUser(), ip, util.DegitalOcean.SSHPort(),
		&ssh.Auth{Keys: util.DefaultSSHPrivateKeys()})

	if err != nil {
		return err
	}

	influxRun := []string{
		"docker", "run", "-d",
		"-p", "8083:8083",
		"-p", "8086:8086",
		"--expose", "8090",
		"--expose", "8099",
		"-e", "PRE_CREATE_DB=telegraf",
		"--name", "influxsrv",
		"tutum/influxdb", // TODO fix version
	}

	_, err = sshcli.Output(strings.Join(influxRun, " "))
	if err != nil {
		return err
	}
	fmt.Println("Starting dashboard database ...")

	grafanaRun := []string{
		"docker", "run", "-d",
		"-p", "80:3000",
		"-e", "HTTP_USER=admin",
		"-e", "HTTP_PASS=admin",
		"-e", "INFLUXDB_HOST=" + ip,
		"-e", "INFLUXDB_PORT=8086",
		"-e", "INFLUXDB_NAME=telegraf",
		"-e", "INFLUXDB_USER=root",
		"-e", "INFLUXDB_PASS=root",
		"--name", "grafana",
		"grafana/grafana", // TODO fix version
	}

	_, err = sshcli.Output(strings.Join(grafanaRun, " "))
	fmt.Println("Starting dashboard user interface ...")

	fmt.Print("Waiting for dashboard service to be ready ...")
	for {
		result, _ := sshcli.Output("docker inspect -f={{.State.Status}} influxsrv")
		result2, _ := sshcli.Output("docker inspect -f={{.State.Status}} grafana")
		if strings.TrimSpace(result) == "running" && strings.TrimSpace(result2) == "running" {
			break
		}
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
	fmt.Println()

	// define database
	tmpl, err := template.New("defineDatasource").Parse(dataSource)
	json := bytes.NewBufferString("")
	err = tmpl.Execute(json, map[string]string{"host": ip})
	err = DefineDB(ip, json.Bytes())
	if err == nil {
		fmt.Println("Datasource defined ...")
	}

	board, err := fixDashboardJson()
	if err != nil {
		return err
	}

	err = DefineDashboard(ip, board)
	if err == nil {
		fmt.Println("Dashboard defined ...")
	}

	return nil
}

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "setup a monitoring dashboard",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// belt dashboard mon
		node := args[0]
		err := InstallDashboard(node)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Dashboard is now running at http://" + GetIP(node) + "/dashboard/db/belt")
	},
}

func DefineDB(ip string, json []byte) error {
	client := &http.Client{}
	req, err := buildRequest("POST", "http://"+ip+"/api/datasources", bytes.NewReader(json))
	if err != nil {
		return err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	_, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	// OK
	if rsp.StatusCode >= 200 && rsp.StatusCode <= 299 {
		return nil
	}

	return fmt.Errorf("NYI")
}

func DefineDashboard(ip string, json []byte) error {
	client := &http.Client{}
	req, err := buildRequest("POST", "http://"+ip+"/api/dashboards/db", bytes.NewReader(json))
	if err != nil {
		return err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	_, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	// OK
	if rsp.StatusCode >= 200 && rsp.StatusCode <= 299 {
		return nil
	}

	return fmt.Errorf("NYI")
}

func buildRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin")

	return req, nil
}

func init() {
	RootCmd.AddCommand(dashboardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dashboardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dashboardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
