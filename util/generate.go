package util

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Env struct{}

var DegitalOcean Env

func (e Env) get(key, def string) string {
	env := viper.GetString(key)
	if env == "" {
		return def
	}
	return env
}

func (e Env) AccessToken() string {
	return e.get("digitalocean.access-token", "")
}

func (e Env) Region() string {
	return e.get("digitalocean.region", "nyc3")
}

func (e Env) Image() string {
	return e.get("digitalocean.image", "ubuntu-15-10-x64")
}

func (e Env) Size() string {
	return e.get("digitalocean.size", "512mb")
}

func (e Env) SSHUser() string {
	return e.get("digitalocean.ssh_user", "root")
}

func (e Env) SSHPort() int {
	result, err := strconv.Atoi(e.get("digitalocean.ssh_port", "22"))
	if err != nil {
		return -1
	}
	return result
}

func (e Env) SSHKey() string {
	return e.get("digitalocean.ssh_key_fingerprint", "")
}

// --digitalocean-ipv6	DIGITALOCEAN_IPV6	false
// --digitalocean-private-networking	DIGITALOCEAN_PRIVATE_NETWORKING	false
// --digitalocean-backups	DIGITALOCEAN_BACKUPS	false

// Generate takes care of IP generation
func Generate(pattern string) []string {
	// fmt.Println("pattern = " + pattern)
	re, _ := regexp.Compile(`\[(.+):(.+)\]`)
	submatch := re.FindStringSubmatch(pattern)
	if submatch == nil {
		return []string{pattern}
	}

	from, err := strconv.Atoi(submatch[1])
	if err != nil {
		return []string{pattern}
	}
	to, err := strconv.Atoi(submatch[2])
	if err != nil {
		return []string{pattern}
	}

	template := re.ReplaceAllString(pattern, "%d")

	var result []string
	for val := from; val <= to; val++ {
		entry := fmt.Sprintf(template, val)
		result = append(result, entry)
	}

	return result
}

const ACTIVE_HOST_FILE = ".belt/active"

func get(file, key string) (string, error) {
	envs := make(map[string]string)
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	err = yaml.Unmarshal(bytes, &envs)
	if err == nil {
		return envs[key], nil
	}

	return "", err
}

func set(file string, key string, value string) error {
	envs := make(map[string]string)
	bytes, err := ioutil.ReadFile(file)
	if err == nil {
		err = yaml.Unmarshal(bytes, &envs)
		if err != nil {
			return err
		}
	}

	envs[key] = value

	data, err := yaml.Marshal(envs)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, 0644)
}

func GetActiveCluster() (string, error) {
	node, err := get(ACTIVE_HOST_FILE, "cluster")
	if err != nil || strings.TrimSpace(node) == "" {
		return "", errors.New("There is no active cluster.")
	}

	return node, nil
}

func GetActive() (string, error) {
	cluster, err := GetActiveCluster()
	if err != nil {
		return "", err
	}
	node, err := get(".belt/"+cluster+"/active", "host")
	if err != nil || strings.TrimSpace(node) == "" {
		return "", fmt.Errorf("%s: there is no active node.", cluster)
	}

	return node, nil
}

func GetActiveByCluster(cluster string) (string, error) {
	node, err := get(".belt/"+cluster+"/active", "host")
	if err != nil || strings.TrimSpace(node) == "" {
		return "", fmt.Errorf("%s: there is no active node.", cluster)
	}

	return node, nil
}

func SetActive(node string) error {
	cluster, err := GetActiveCluster()
	if err != nil {
		return err
	}
	return set(".belt/"+cluster+"/active", "host", node)
}

func SetActiveCluster(cluster string) error {
	return set(ACTIVE_HOST_FILE, "cluster", cluster)
}

func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		if os.Getenv("CYGWIN") != ""  || os.Getenv("TERM") == "cygwin" {
			bout, err := exec.Command("cygpath", "-w", os.Getenv("HOME")).Output()
			if err != nil {
				return ""
			}
			return strings.TrimSpace(string(bout))
		}

		if os.Getenv("HOME") != "" {
			return os.Getenv("HOME")
		}

		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func DefaultSSHPrivateKeys() []string {
	return []string{path.Join(GetHomeDir(), ".ssh", "id_rsa")}
}
