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

type Repo struct {
	Remote string
	Dir    string
}

func findGitFolder(path string) (string, error) {
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i > 0; i-- {
		thisPath := "/" + filepath.Join(parts[0:i]...)
		thisPathGit := filepath.Join(thisPath, ".git")
		if info, err := os.Stat(thisPathGit); err == nil && info.IsDir() {
			return thisPath, nil
		}
	}

	return "", errors.New("Could not find .git folder")
}

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

func (r *Repo) URL(file string) (string, error) {
	absFile, err := filepath.Abs(file)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"file": file,
		}).Panic("Cannot calculate abs path")

	}

	relativeFile := strings.TrimPrefix(strings.Replace(absFile, r.Dir, "", -1), "/")

	gl := gitlab.Remote{R: r.Remote}
	res, err := gl.File(r.Branch(), relativeFile)
	if err == nil {
		return res, nil
	}

	cc := codecommit.Remote{R: r.Remote}
	res, err = cc.File(r.Branch(), relativeFile)
	if err == nil {
		return res, nil
	}

	gh := github.Remote{R: r.Remote}
	res, err = gh.File(r.Branch(), relativeFile)
	if err == nil {
		return res, nil
	}

	return "", errors.New("Cannot handle remote type")
}
