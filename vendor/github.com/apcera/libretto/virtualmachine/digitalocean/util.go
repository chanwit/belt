// Copyright 2015 Apcera Inc. All rights reserved.

package digitalocean

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// BuildRequest builds an http request for this provider.
func BuildRequest(token, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}

// Update vm.Droplet values. This occurs in GetState(), so we call that and
// ignore the state string.
func (vm *VM) Update() error {
	_, err := vm.GetState()
	return err
}

// GetDroplet returns a single droplet
func GetDroplet(token, id string) (*Droplet, error) {
	client := &http.Client{}
	req, err := BuildRequest(token, "GET", apiBaseURL+apiDropletURL+"/"+id, nil)
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
	if rsp.Status[0] != StatusOk {
		return nil, fmt.Errorf("Error: %s: %s", rsp.Status, string(b))
	}

	r := &DropletResponse{}
	err = json.Unmarshal(b, r)
	if err != nil {
		return nil, err
	}
	return r.Droplet, nil
}

// GetDroplets returns and array of droplets
func GetDroplets(token string) (*DropletsResponse, error) {
	client := &http.Client{}
	req, err := BuildRequest(token, "GET", apiBaseURL+apiDropletURL, nil)
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
	if rsp.Status[0] != StatusOk {
		return nil, fmt.Errorf("Error: %s: %s", rsp.Status, string(b))
	}

	r := &DropletsResponse{}
	err = json.Unmarshal(b, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PrintDroplet prints the basic droplet values
func PrintDroplet(droplet *Droplet) {
	fmt.Println("ID:", droplet.ID)
	fmt.Println("Name:", droplet.Name)
	fmt.Println("Status:", droplet.Status)
	fmt.Println("Locked:", fmt.Sprintf("%t", droplet.Locked))
	fmt.Println("CreatedAt:", droplet.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("SizeSlug:", droplet.Size.Slug)
	fmt.Println("Region:", droplet.Region.Name)
	fmt.Println("Image:", droplet.Image.Name)
	for i, ip := range droplet.Networks.V4 {
		fmt.Println(fmt.Sprintf("IP: %d", i+1), ip.IPAddress, ip.Type)
	}
	for i, ip := range droplet.Networks.V6 {
		fmt.Println(fmt.Sprintf("IP: %d", i+1), ip.IPAddress, ip.Type)
	}
}
