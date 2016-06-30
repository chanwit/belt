package util

import (
	"fmt"
	"os"

	"regexp"
	"strconv"
)

type Env struct {}

var DefaultEnv Env

func (e Env) get(key, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}
	return env
}

func (e Env) AccessToken() string {
	return e.get("DIGITALOCEAN_ACCESS_TOKEN", "")
}

func (e Env) Region() string {
	return e.get("DIGITALOCEAN_REGION", "nyc3")
}

func (e Env) Image() string {
	return e.get("DIGITALOCEAN_IMAGE", "ubuntu-15-10-x64")
}

func (e Env) Size() string {
	return e.get("DIGITALOCEAN_SIZE", "512mb")
}

func (e Env) SSHUser() string {
	return e.get("DIGITALOCEAN_SSH_USER", "root")
}

func (e Env) SSHPort() string {
	return e.get("DIGITALOCEAN_SSH_PORT", "22")
}

func (e Env) SSHKey() string {
	return e.get("DIGITALOCEAN_SSH_KEY_FINGERPRINT", "")
}

// --digitalocean-ipv6	DIGITALOCEAN_IPV6	false
// --digitalocean-private-networking	DIGITALOCEAN_PRIVATE_NETWORKING	false
// --digitalocean-backups	DIGITALOCEAN_BACKUPS	false

// Generate takes care of IP generation
func Generate(pattern string) []string {
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
