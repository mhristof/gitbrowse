package codecommit

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/mhristof/gitbrowse/log"
)

var (
	// ErrorNotCodeCommit Error returned when the remote is not a codecommit remote
	ErrorNotCodeCommit = errors.New("not a codecommit remote")
)

// Remote Holds remote information
type Remote struct {
	R string
}

// Valid Checks if the remote is avalid codecommit remote
func (r *Remote) Valid() bool {
	parts, err := url.Parse(r.R)
	if err != nil {
		return false
	}

	if strings.HasPrefix(parts.Host, "git-codecommit") && strings.HasSuffix(parts.Host, ".amazonaws.com") {
		return true
	}

	return false
}

// Region Retrieves the region of the remote
func (r *Remote) Region() string {
	parts, err := url.Parse(r.R)
	if err != nil {
		log.WithFields(log.Fields{
			"r.R": r.R,
		}).Error("Cannot extract region")

	}

	return strings.Split(parts.Host, ".")[1]
}

// URL Sanitise the remote to a valid url
func (r *Remote) URL() string {
	var remRegex = regexp.MustCompile(`https://git-codecommit.(?P<region>.*).amazonaws.com/v1/repos/(?P<repo>.*)`)
	match := remRegex.FindStringSubmatch(r.R)

	info := map[string]string{}

	if remRegex.MatchString(r.R) {
		for i, name := range remRegex.SubexpNames() {
			info[name] = match[i]
		}
	}

	if region, _ := info["region"]; region == "" {
		log.WithFields(log.Fields{
			"r.R": r.R,
		}).Error("Cannot find region")
	}

	if repo, _ := info["repo"]; repo == "" {
		log.WithFields(log.Fields{
			"r.R": r.R,
		}).Error("Cannot retrieve repo")

	}

	return fmt.Sprintf("https://%s.console.aws.amazon.com/codesuite/codecommit/repositories/%s", info["region"], info["repo"])
}

// File Retrieves the URL for the given file/branch combination
func (r *Remote) File(branch, file string, line int) (string, error) {
	if !r.Valid() {
		return "", ErrorNotCodeCommit
	}

	branch = strings.Replace(branch, "refs/heads/", "", -1)
	ret := fmt.Sprintf("%s/browse/refs/heads/%s/--/%s?region=%s", r.URL(), branch, file, r.Region())

	if line >= 0 {
		ret += fmt.Sprintf("#L%d-%d", line, line)
	}

	return ret, nil
}
