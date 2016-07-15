package drivers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/apcera/libretto/virtualmachine/digitalocean"
)

// Base API URL strings
const (
	apiBaseURL    = "https://api.digitalocean.com"
	apiDropletURL = "v2/droplets"
)

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