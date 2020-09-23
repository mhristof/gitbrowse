package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mhristof/germ/log"
	"gopkg.in/ini.v1"
)

type Repo struct {
	Remote string
	Dir    string
}

func findGitFolder(path string) (string, error) {
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i > 0; i-- {
		thisPath := "/" + filepath.Join(parts[0:i]...)
		thisPathGit := filepath.Join(thisPath, ".git")
		fmt.Println(thisPathGit)
		if info, err := os.Stat(thisPathGit); err == nil && info.IsDir() {
			return thisPath, nil
		}
	}

	return "", errors.New("Could not find .git folder")
}

func New(directory string) (*Repo, error) {
	absDir, err := filepath.Abs(directory)
	if err != nil {
		panic(err)
	}

	if info, err := os.Stat(absDir); err != nil || !info.IsDir() {
		absDir, err = findGitFolder(absDir)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(absDir)

	cfg, err := ini.Load(filepath.Join(absDir, ".git/config"))
	if err != nil {
		log.WithFields(log.Fields{
			"absDir": absDir,
		}).Error("Cant read .git/config")
	}

	return &Repo{
		Remote: cfg.Section(`remote "origin"`).Key("url").Value(),
		Dir:    absDir,
	}, nil
}

func (r *Repo) Branch() string {
	head, err := ioutil.ReadFile(filepath.Join(r.Dir, ".git/HEAD"))
	if err != nil {
		panic(err)
	}

	headS := strings.Split(strings.Split(string(head), "\n")[0], " ")[1]
	return headS
}

func (r *Repo) URL(file string) (string, error) {

	res, err := r.gitlab(file)
	if err == nil {
		return res, nil
	}

	return "", errors.New("Cannot handle remote type")
}

func (r *Repo) gitlab(file string) (string, error) {
	if !isGitlab(r.Remote) {
		return "", errors.New("Not a gitlab remote")
	}

	absFile, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}

	relativeFile := strings.TrimPrefix(strings.Replace(absFile, r.Dir, "", -1), "/")
	fmt.Println("relative file", relativeFile)
	remote := gitlabRemoteToURL(r.Remote)

	branch := strings.Replace(r.Branch(), "refs/heads/origin/", "", -1)
	return fmt.Sprintf("%s/-/blob/%s/%s", remote, branch, relativeFile), nil
}

func gitlabRemoteToURL(remote string) string {
	var remRegex = regexp.MustCompile(`https://(?P<username>.*):(?P<token>.*)@(?P<url>.*)`)
	match := remRegex.FindStringSubmatch(remote)

	if remRegex.MatchString(remote) {
		for i, name := range remRegex.SubexpNames() {
			if name == "url" {
				return fmt.Sprintf("https://%s", match[i])
			}

		}
	}

	panic("panic!")

}

func isGitlab(remote string) bool {
	var gitlabURL = regexp.MustCompile(`gitlab`)

	if gitlabURL.MatchString(remote) {
		return true
	}

	return false
}
