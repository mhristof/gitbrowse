package github

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mhristof/gitbrowse/log"
)

// ErrorNotGithub Error thrown when the remote doesnt seem like a valid github remote
var ErrorNotGithub = errors.New("Not a valid github remote")

// Remote Holds remote information
type Remote struct {
	R string
}

// Valid Checks if the remote is a valid Github repo
func (r *Remote) Valid() bool {
	url := regexp.MustCompile(`github`)

	if url.MatchString(r.R) {
		return true
	}

	return false
}

func gitToHTTP(remote string) string {
	if !strings.HasPrefix(remote, "git@") {
		return remote
	}

	user := strings.Split(filepath.Dir(remote), ":")[1]
	host := strings.Split(strings.Split(remote, ":")[0], "@")[1]
	repo := filepath.Base(remote)

	return fmt.Sprintf("https://%s/%s/%s", host, user, repo)
}

// URL Sanitises the remote to a valid url
func (r *Remote) URL() string {
	remote := gitToHTTP(r.R)
	remRegex := regexp.MustCompile(`https://(?P<url>.*)`)
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

// File Retrieve the URL for the given file/branch combination
func (r *Remote) File(branch, file string, line int) (string, error) {
	if !r.Valid() {
		return "", ErrorNotGithub
	}

	branch = strings.Replace(branch, "refs/heads/", "", -1)
	ret := fmt.Sprintf("%s/tree/%s/%s", r.URL(), branch, file)
	if line >= 0 {
		ret += fmt.Sprintf("#L%d", line)
	}

	log.WithFields(log.Fields{
		"branch": branch,
		"file":   file,
		"line":   line,
		"remote": r,
	}).Debug("calculated url")
	return ret, nil
}
