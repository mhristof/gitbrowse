package git

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mhristof/gitbrowse/codecommit"
	"github.com/mhristof/gitbrowse/github"
	"github.com/mhristof/gitbrowse/gitlab"
	"github.com/mhristof/gitbrowse/log"
	"gopkg.in/ini.v1"
)

// Repo holds information about a repository
type Repo struct {
	Remote string
	Dir    string
}

var (
	// ErrorNotAGitRepo is thrown when the given folder/config is not a git repository
	ErrorNotAGitRepo = errors.New("not a git repository")
)

func findGitFolder(path string) (string, error) {
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i > 0; i-- {
		thisPath := "/" + filepath.Join(parts[0:i]...)
		thisPathGit := filepath.Join(thisPath, ".git")
		if info, err := os.Stat(thisPathGit); err == nil && info.IsDir() {
			return thisPath, nil
		}
	}

	return "", ErrorNotAGitRepo
}

// New Create a new git repository object from the given directory.
// The directory could be relative or absolute folder or file inside the git
// repository
func New(directory string) (*Repo, error) {
	absDir, err := filepath.Abs(directory)
	if err != nil {
		log.WithFields(log.Fields{
			"err":       err,
			"directory": directory,
		}).Panic("Cannot calculate abs path")

	}

	if info, err := os.Stat(absDir); err != nil || !info.IsDir() {
		absDir, err = findGitFolder(absDir)
		if err != nil {
			log.WithFields(log.Fields{
				"err":    err,
				"absDir": absDir,
			}).Panic("Cannot find .git folder")

		}
	}

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

// Branch Return the current branch of the git repository by reading .git/HEAD
func (r *Repo) Branch() string {
	head, err := ioutil.ReadFile(filepath.Join(r.Dir, ".git/HEAD"))
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"r.Dir": r.Dir,
		}).Panic("Cannot find .git/HEAD")

	}

	headS := strings.Split(strings.Split(string(head), "\n")[0], " ")[1]
	return headS
}

// URL Returns the web url for the given file. Currently gitlab, github and codecommit
// are supported
func (r *Repo) URL(file string, line int) (string, error) {
	absFile, err := filepath.Abs(file)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"file": file,
		}).Panic("Cannot calculate abs path")

	}

	relativeFile := strings.TrimPrefix(strings.Replace(absFile, r.Dir, "", -1), "/")

	gl := gitlab.Remote{R: r.Remote}
	res, err := gl.File(r.Branch(), relativeFile, line)
	if err == nil {
		return res, nil
	}

	cc := codecommit.Remote{R: r.Remote}
	res, err = cc.File(r.Branch(), relativeFile, line)
	if err == nil {
		return res, nil
	}

	gh := github.Remote{R: r.Remote}
	res, err = gh.File(r.Branch(), relativeFile, line)
	if err == nil {
		return res, nil
	}

	return "", errors.New("Cannot handle remote type")
}
