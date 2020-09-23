package gitlab

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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
	var gitlabURL = regexp.MustCompile(`gitlab`)

	if gitlabURL.MatchString(r.R) {
		return true
	}

	return false
}

func (r *Remote) URL() string {
	var remRegex = regexp.MustCompile(`https://(?P<username>.*):(?P<token>.*)@(?P<url>.*)`)
	match := remRegex.FindStringSubmatch(r.R)

	fmt.Println(r.R)
	if remRegex.MatchString(r.R) {
		for i, name := range remRegex.SubexpNames() {
			if name == "url" {
				return fmt.Sprintf("https://%s", match[i])
			}

		}
	}
	panic("Not a gitlab remote")
}

func (r *Remote) File(branch, file string) (string, error) {
	if !r.Valid() {
		return "", errors.New("cannot handle this remote")
	}

	branch = strings.Replace(branch, "refs/heads/origin/", "", -1)
	return fmt.Sprintf("%s/-/blob/%s/%s", r.URL(), branch, file), nil
}
