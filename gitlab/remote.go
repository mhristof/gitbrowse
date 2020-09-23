package gitlab

import (
	"fmt"
	"regexp"
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

func (r *Remote) isGitlab() bool {
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
