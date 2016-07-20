// Copyright 2015 Apcera Inc. All rights reserved.

package digitalocean

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	libssh "github.com/apcera/libretto/ssh"
	"github.com/apcera/libretto/util"
	lvm "github.com/apcera/libretto/virtualmachine"
)

var (
	// ErrNoInstanceID is returned when attempting to perform an operation on an instance, but the ID is missing.
	ErrNoInstanceID = errors.New("Missing droplet ID")
)

// Base API URL strings
const (
	apiBaseURL    = "https://api.digitalocean.com"
	apiDropletURL = "/v2/droplets"
)

// VM struct represents a full DigitalOcean VM in libretto. It contains the
// droplet itself, along with authentication and SSH credential information. It
// also contains the original creation config.
type VM struct {
	APIToken    string // required
	Credentials libssh.Credentials
	Config      Config
	Droplet     *Droplet
}

var _ lvm.VirtualMachine = (*VM)(nil)

// Config is the new droplet payload
type Config struct {
	Name              string   `json:"name,omitempty"`   // required
	Region            string   `json:"region,omitempty"` // required
	Size              string   `json:"size,omitempty"`   // required
	Image             string   `json:"image,omitempty"`  // required
	SSHKeys           []string `json:"ssh_keys,omitempty"`
	Backups           bool     `json:"backups,omitempty"`
	IPv6              bool     `json:"ipv6,omitempty"`
	PrivateNetworking bool     `json:"private_networking,omitempty"`
	UserData          string   `json:"user_data,omitempty"`
}

// DropletsResponse is the API response containing multiple droplets
type DropletsResponse struct {
	Droplets []*Droplet `json:"droplets"`
	Meta     *Meta      `json:"meta,omitempty"`
}

// Meta contains metadata from API responses
type Meta struct {
	Total int `json:"total,omitempty"`
}

// DropletResponse is the API response containing one droplet
type DropletResponse struct {
	Droplet *Droplet `json:"droplet,omitempty"`
}

// Droplet is the base VM type in DigitalOcean
type Droplet struct {
	ID          int       `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Region      *Region   `json:"region,omitempty"`
	Image       *Image    `json:"image,omitempty"`
	Size        *Size     `json:"size,omitempty"`
	SizeSlug    string    `json:"size_slug,omitempty"`
	Locked      bool      `json:"locked,omitempty"`
	Status      string    `json:"status,omitempty"`
	Networks    *Networks `json:"networks,omitempty"`
	Kernel      *Kernel   `json:"kernel,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	BackupIds   []int64   `json:"backup_ids,omitempty"`
	SnapshotIds []int64   `json:"snapshot_ids,omitempty"`
	ActionIds   []int64   `json:"action_ids,omitempty"`
}

// Networks contains information on all droplet networks
type Networks struct {
	V4 []*V4Network `json:"v4,omitempty"`
	V6 []*V6Network `json:"v6,omitempty"`
}

// Kernel contains droplet kernel information
type Kernel struct {
	ID      int64  `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// V4Network contains an IPV4 network's values
type V4Network struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   string `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}

// V6Network contains an IPV6 network's values
type V6Network struct {
	IPAddress string `json:"ip_address,omitempty"`
	Netmask   int64  `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	Type      string `json:"type,omitempty"`
}

// Region contains droplet region information
type Region struct {
	Slug      string   `json:"slug,omitempty"`
	Name      string   `json:"name,omitempty"`
	Sizes     []string `json:"sizes,omitempty"`
	Available bool     `json:"available,omitempty"`
	Features  []string `json:"features,omitempty"`
}

// Image contains droplet image information
type Image struct {
	ID           int       `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Distribution string    `json:"distribution,omitempty"`
	Slug         string    `json:"slug,omitempty"`
	Public       bool      `json:"public,omitempty"`
	Regions      []string  `json:"regions,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

// Size describes droplet sizing
type Size struct {
	Slug          string      `json:"slug,omitempty"`
	Memory        int         `json:"memory,omitempty"`
	VCpus         int         `json:"v_cpus,omitempty"`
	Disk          int         `json:"disk,omitempty"`
	Transfer      interface{} `json:"transfer,omitempty"`
	PriceMonthley float64     `json:"price_monthley,omitempty"`
	PriceHourly   float64     `json:"price_hourly,omitempty"`
	Regions       []string    `json:"regions,omitempty"`
}

const (
	// StatusOk is the first number of a successful HTTP return code
	StatusOk = '2'
	// StatusNotFound is the status code for a 404
	StatusNotFound = 404
)

// GetName returns the name of the virtual machine
func (vm *VM) GetName() string {
	return vm.Config.Name
}

// Provision creates a new VM
func (vm *VM) Provision() error {
	b, err := json.Marshal(vm.Config)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := BuildRequest(vm.APIToken, "POST", apiBaseURL+apiDropletURL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	b, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.Status[0] != StatusOk {
		return fmt.Errorf("Error: %s: %s", rsp.Status, string(b))
	}

	// Fill out vm.Droplet with data on new droplet
	r := &DropletResponse{}
	err = json.Unmarshal(b, r)
	if err != nil {
		return err
	}
	vm.Droplet = r.Droplet
	return nil
}

// GetIPs returns a list of ip addresses associated with the VM
func (vm *VM) GetIPs() ([]net.IP, error) {
	var ips []net.IP
	if err := vm.Update(); err != nil {
		return nil, err
	}
	for _, ip := range vm.Droplet.Networks.V4 {
		ips = append(ips, net.ParseIP(ip.IPAddress))
	}
	for _, ip := range vm.Droplet.Networks.V6 {
		ips = append(ips, net.ParseIP(ip.IPAddress))
	}
	return ips, nil
}

// GetSSH returns an ssh client for the the vm.
func (vm *VM) GetSSH(options libssh.Options) (libssh.Client, error) {
	ips, err := util.GetVMIPs(vm, options)
	if err != nil {
		return nil, err
	}

	client := libssh.SSHClient{Creds: &vm.Credentials, IP: ips[0], Port: 22, Options: options}
	return &client, nil
}

// Destroy powers off the VM and deletes its files from disk
func (vm *VM) Destroy() error {
	id := fmt.Sprintf("%v", vm.Droplet.ID)
	if id == "" {
		return ErrNoInstanceID
	}

	client := &http.Client{}
	req, err := BuildRequest(vm.APIToken, "DELETE", apiBaseURL+apiDropletURL+"/"+id, nil)
	if err != nil {
		return err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.Status[0] != StatusOk {
		return fmt.Errorf("Error: %s: %s", rsp.Status, string(b))
	}

	return nil
}

// GetState gets the running state of the VM through the DigitalOcean API
// Returns droplet state if available and 'not_found' if ID could not be located.
func (vm *VM) GetState() (string, error) {
	id := fmt.Sprintf("%v", vm.Droplet.ID)
	if id == "" {
		return "", ErrNoInstanceID
	}

	client := &http.Client{}
	req, err := BuildRequest(vm.APIToken, "GET", apiBaseURL+apiDropletURL+"/"+id, nil)
	if err != nil {
		return "", err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	if rsp.StatusCode == StatusNotFound {
		return "not_found", nil
	}
	if rsp.Status[0] != StatusOk {
		return "", fmt.Errorf("Error: %s: %s", rsp.Status, string(b))
	}

	// Fill out vm.Droplet with data on droplet
	r := &DropletResponse{}
	err = json.Unmarshal(b, r)
	if err != nil {
		return "", err
	}
	vm.Droplet = r.Droplet
	return vm.Droplet.Status, nil
}

// Start powers on the VM
func (vm *VM) Start() error {
	id := fmt.Sprintf("%v", vm.Droplet.ID)
	if id == "" {
		return ErrNoInstanceID
	}

	client := &http.Client{}
	req, err := BuildRequest(vm.APIToken, "POST", apiBaseURL+apiDropletURL+"/"+id+"/actions", strings.NewReader(`{"type": "power_on"}`))
	if err != nil {
		return err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.Status[0] != StatusOk {
		return fmt.Errorf("Error: %s: %s", rsp.Status, string(b))
	}

	return nil
}

// Halt powers off the VM without destroying it
func (vm *VM) Halt() error {
	id := fmt.Sprintf("%v", vm.Droplet.ID)
	if id == "" {
		return ErrNoInstanceID
	}

	client := &http.Client{}
	req, err := BuildRequest(vm.APIToken, "POST", apiBaseURL+apiDropletURL+"/"+id+"/actions", strings.NewReader(`{"type": "power_off"}`))
	if err != nil {
		return err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if rsp.Status[0] != StatusOk {
		return fmt.Errorf("Error: %s: %s", rsp.Status, string(b))
	}

	return nil
}

// Suspend always returns an error because this isn't supported by DigitalOcean
func (vm *VM) Suspend() error {
	return lvm.ErrSuspendNotSupported
}

// Resume always returns an error because this isn't supported by DigitalOcean
func (vm *VM) Resume() error {
	return lvm.ErrResumeNotSupported
}
