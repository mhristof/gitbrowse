package github

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mhristof/germ/log"
)

type Remote struct {
	R string
}

func New(in string) Remote {
	return Remote{
		R: in,
	}
}
func (r Remote) String() string {
	return r.R
}

func (r *Remote) Valid() bool {
	var url = regexp.MustCompile(`github`)

	if url.MatchString(r.R) {
		return true
	}

	return false
}

func gitToHttp(remote string) string {
	if !strings.HasPrefix(remote, "git@") {
		return remote
	}

	user := strings.Split(filepath.Dir(remote), ":")[1]
	host := strings.Split(strings.Split(remote, ":")[0], "@")[1]
	repo := filepath.Base(remote)

	return fmt.Sprintf("https://%s/%s/%s", host, user, repo)
}

func (r *Remote) URL() string {
	remote := gitToHttp(r.R)
	var remRegex = regexp.MustCompile(`https://(?P<url>.*)`)
	match := remRegex.FindStringSubmatch(remote)

	if remRegex.MatchString(remote) {
		for i, name := range remRegex.SubexpNames() {
			if name == "url" {
				return strings.Replace(fmt.Sprintf("https://%s", match[i]), ".git", "", -1)
			}

		}
	}

	log.WithFields(log.Fields{
		"r.R": r.R,
	}).Error("Not a github remote")

	return ""
}

func (r *Remote) File(branch, file string) (string, error) {
	if !r.Valid() {
		return "", errors.New("cannot handle this remote")
	}

	branch = strings.Replace(branch, "refs/heads/", "", -1)
	return fmt.Sprintf("%s/blob/%s/%s", r.URL(), branch, file), nil
}
