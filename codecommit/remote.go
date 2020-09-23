package codecommit

import (
	"errors"
	"fmt"
	"net/url"
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

func (r *Remote) Valid() bool {
	parts, err := url.Parse(r.R)
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(parts.Host, "git-codecommit") && strings.HasSuffix(parts.Host, ".amazonaws.com") {
		return true
	}

	return false
}

func (r *Remote) Region() string {
	parts, err := url.Parse(r.R)
	if err != nil {
		panic(err)
	}

	return strings.Split(parts.Host, ".")[1]
}

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
		panic("Not a codecommit remote - missing region")
	}

	if repo, _ := info["repo"]; repo == "" {
		panic("Not a codecommit remote - missing repo")
	}

	return fmt.Sprintf("https://%s.console.aws.amazon.com/codesuite/codecommit/repositories/%s", info["region"], info["repo"])
}

func (r *Remote) File(branch, file string) (string, error) {
	if !r.Valid() {
		return "", errors.New("cannot handle this remote")
	}

	//fmt.Println(fmt.Sprintf("branch: %+v", branch))

	branch = strings.Replace(branch, "refs/heads/", "", -1)
	return fmt.Sprintf("%s/browse/refs/heads/%s/--/%s?region=%s", r.URL(), branch, file, r.Region()), nil
}
