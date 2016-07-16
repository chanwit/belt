package drivers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"

	"github.com/apcera/libretto/virtualmachine/digitalocean"
)

// Base API URL strings
const (
	apiBaseURL    = "https://api.digitalocean.com"
	apiDropletURL = "v2/droplets"
)

type Config struct {
	Names             []string   `json:"names,omitempty"` // required
	Region            string   `json:"region,omitempty"`  // required
	Size              string   `json:"size,omitempty"`    // required
	Image             string   `json:"image,omitempty"`   // required
	SSHKeys           []string `json:"ssh_keys,omitempty"`
	Backups           bool     `json:"backups,omitempty"`
	IPv6              bool     `json:"ipv6,omitempty"`
	PrivateNetworking bool     `json:"private_networking,omitempty"`
	UserData          string   `json:"user_data,omitempty"`
}

func Provision(token string, config Config) (*digitalocean.DropletsResponse, error) {
	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := digitalocean.BuildRequest(token, "POST", apiBaseURL + "/" + apiDropletURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	b, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if rsp.Status[0] != digitalocean.StatusOk {
		return nil, fmt.Errorf("Error: %s: %s", rsp.Status, string(b))
	}

	// Fill out vm.Droplet with data on new droplet
	r := &digitalocean.DropletsResponse{}
	err = json.Unmarshal(b, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func GetTotalDroplets(token string) (int, error) {
	client := &http.Client{}
	req, err := digitalocean.BuildRequest(token, "GET", apiBaseURL+ "/" + apiDropletURL + "?page=1&per_page=1", nil)
	if err != nil {
		return 0, err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return 0, err
	}
	if rsp.Status[0] != digitalocean.StatusOk {
		return 0, fmt.Errorf("Error: %s: %s", rsp.Status, string(b))
	}

	r := &digitalocean.DropletsResponse{}
	err = json.Unmarshal(b, r)
	if err != nil {
		return 0, err
	}

	return r.Meta.Total, nil
}

func GetAllDroplets(token string) (*digitalocean.DropletsResponse, error) {
	total, err := GetTotalDroplets(token)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return nil, nil
	}

	client := &http.Client{}
	req, err := digitalocean.BuildRequest(token, "GET", apiBaseURL+ "/" + apiDropletURL + fmt.Sprintf("?page=1&per_page=%d", total), nil)
	if err != nil {
		return nil, err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if rsp.Status[0] != digitalocean.StatusOk {
		return nil, fmt.Errorf("Error: %s: %s", rsp.Status, string(b))
	}

	r := &digitalocean.DropletsResponse{}
	err = json.Unmarshal(b, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

