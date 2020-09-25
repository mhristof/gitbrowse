package git

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mhristof/gitbrowse/codecommit"
	"github.com/mhristof/gitbrowse/github"
	"github.com/mhristof/gitbrowse/gitlab"
	"github.com/mhristof/gitbrowse/log"
	"github.com/pkg/errors"
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
	for i := len(parts); i > 0; i-- {
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
		return nil, err
	}

	absDir, err = findGitFolder(absDir)
	if err != nil {
		return nil, errors.Wrap(err, "Canot find .git folder in "+directory)
	}

	cfg, err := ini.Load(filepath.Join(absDir, ".git/config"))
	if err != nil {
		return nil, errors.Wrap(err, "Cant read .git/config")
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

type Remote interface {
	File(string, string, int) (string, error)
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

	remotes := []Remote{
		&gitlab.Remote{R: r.Remote},
		&codecommit.Remote{R: r.Remote},
		&github.Remote{R: r.Remote},
	}

	var wg sync.WaitGroup
	wg.Add(len(remotes))
	res := make(chan string, 1)

	for _, remote := range remotes {
		go url(&wg, r, remote, relativeFile, line, res)
	}
	wg.Wait()

	return <-res, nil
}

func url(wg *sync.WaitGroup, r *Repo, remote Remote, relativeFile string, line int, c chan string) {
	defer wg.Done()

	this, err := remote.File(r.Branch(), relativeFile, line)
	if err != nil {
		return
	}
	c <- this
}
