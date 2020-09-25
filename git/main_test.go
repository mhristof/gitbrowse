package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mhristof/gitbrowse/log"
	"github.com/stretchr/testify/assert"
)

func eval(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"args": args,
		}).Panic("Cannot execute command")

	}
}

func newRepo(remote, branch, file string) string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	eval([]string{"bash", "-c", fmt.Sprintf("cd %s && git init", dir)})

	eval([]string{"git", "-C", dir, "remote", "add", "origin", remote})

	eval([]string{"git", "-C", dir, "checkout", "-b", branch})
	eval([]string{"mkdir", "-p", filepath.Join(dir, filepath.Dir(file))})
	eval([]string{"touch", filepath.Join(dir, file)})

	return dir
}

func TestURL(t *testing.T) {
	var cases = []struct {
		name     string
		remote   string
		file     string
		path     string
		branch   string
		expected string
	}{
		{
			name:     "gitlab https remote with username and token on master",
			remote:   "https://glUsernames:glToken@gitlab.com/project/repository",
			branch:   "master",
			file:     "README.md",
			expected: "https://gitlab.com/project/repository/-/blob/master/README.md#L0",
		},
		{
			name:     "gitlab https remote with username and token on a branch",
			remote:   "https://glUsernames:glToken@gitlab.com/project/repository",
			branch:   "branch",
			file:     "README.md",
			expected: "https://gitlab.com/project/repository/-/blob/branch/README.md#L0",
		},
		{
			name:     "github remote",
			remote:   "https://github.com/user/repo.git",
			branch:   "branch",
			file:     "README.md",
			expected: "https://github.com/user/repo/blob/branch/README.md#L0",
		},
		{
			name:     "codecommit repo",
			remote:   "https://git-codecommit.region.amazonaws.com/v1/repos/foobar",
			branch:   "branch",
			file:     "test/readme.md",
			expected: "https://region.console.aws.amazon.com/codesuite/codecommit/repositories/foobar/browse/refs/heads/branch/--/test/readme.md?region=region#L0-0",
		},
	}

	for _, test := range cases {
		dir := newRepo(test.remote, test.branch, test.file)
		defer os.RemoveAll(dir)

		repo, err := New(dir)

		assert.Nil(t, err, test.name)

		url, err := repo.URL(filepath.Join(dir, test.file), 0)
		assert.Nil(t, err, test.name)

		assert.Equal(t, test.expected, url, test.name)
	}
}

func TestFindGitFolder(t *testing.T) {
	var cases = []struct {
		name string
		dir  string
		file string
		err  error
	}{
		{
			name: "file in root",
			dir:  newRepo("", "master", "foo"),
			file: "foo",
			err:  nil,
		},
		{
			name: "file in folder",
			dir:  newRepo("", "master", "foo/bar"),
			file: "foo/bar",
			err:  nil,
		},
		{
			name: "file not in repo",
			dir:  "/tmp",
			file: "/tmp/foobar",
			err:  ErrorNotAGitRepo,
		},
	}

	for _, test := range cases {
		abs, err := filepath.Abs(test.file)
		path, err := findGitFolder(filepath.Join(test.dir, abs))

		assert.Equal(t, test.err, err, test.name)
		if test.err == nil {
			assert.Equal(t, test.dir, path, test.name)
		}
	}
}
