package gitlab

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mhristof/germ/log"
)

var (
	// ErrorNotGitlab The remote doesnt seem like a Gitlab server
	ErrorNotGitlab = errors.New("Not a valid Gitlab remote")
)

// Remote Represent a gitlab remote
type Remote struct {
	R string
}

// Valid Checks a remote to see if its a valid gitlab instance
func (r *Remote) Valid() bool {
	var gitlabURL = regexp.MustCompile(`gitlab`)

	if gitlabURL.MatchString(r.R) {
		return true
	}

	return false
}

// URL Return the URL of the remote by sanitizing it
func (r *Remote) URL() string {
	remote := gitlabHTTP(r.R)
	if remote != "" {
		return remote
	}

	remote = gitlabSSH(r.R)
	if remote != "" {
		return remote
	}

	log.WithFields(log.Fields{
		"r.R": r.R,
	}).Error("Not a gitlab remote")

	return ""

}

func gitlabSSH(url string) string {
	if !strings.HasPrefix(url, "git@gitlab.com:") {
		return ""
	}

	url = strings.TrimSuffix(url, ".git")

	return strings.Replace(url, "git@gitlab.com:", "https://gitlab.com/", 1)
}

func gitlabHTTP(url string) string {
	var remRegex = regexp.MustCompile(`https://(?P<username>.*):(?P<token>.*)@(?P<url>.*)`)
	match := remRegex.FindStringSubmatch(url)

	if remRegex.MatchString(url) {
		for i, name := range remRegex.SubexpNames() {
			if name == "url" {
				return fmt.Sprintf("https://%s", match[i])
			}

		}
	}
	return ""
}

// File Retrieves the file url for the given file. Throws a ErrorNotGitlab
// if the repository is not a valid gitlab url
func (r *Remote) File(branch, file string, line int) (string, error) {
	if !r.Valid() {
		return "", ErrorNotGitlab
	}

	branch = strings.Replace(branch, "refs/heads/", "", -1)
	ret := fmt.Sprintf("%s/-/blob/%s/%s", r.URL(), branch, file)

	if line >= 0 {
		ret += fmt.Sprintf("#L%d", line)
	}

	return ret, nil
}
