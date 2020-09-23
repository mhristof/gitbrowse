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

func newRepo(remote, branch string) string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	eval([]string{"bash", "-c", fmt.Sprintf("cd %s && git init", dir)})

	eval([]string{"git", "-C", dir, "remote", "add", "origin", remote})

	eval([]string{"git", "-C", dir, "checkout", "-b", branch})

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
			expected: "https://gitlab.com/project/repository/-/blob/master/README.md",
		},
		{
			name:     "gitlab https remote with username and token on a branch",
			remote:   "https://glUsernames:glToken@gitlab.com/project/repository",
			branch:   "branch",
			file:     "README.md",
			expected: "https://gitlab.com/project/repository/-/blob/branch/README.md",
		},
		{
			name:     "github remote",
			remote:   "https://github.com/user/repo.git",
			branch:   "branch",
			file:     "README.md",
			expected: "https://github.com/user/repo/blob/branch/README.md",
		},
		{
			name:     "codecommit repo",
			remote:   "https://git-codecommit.region.amazonaws.com/v1/repos/foobar",
			branch:   "branch",
			file:     "test/readme.md",
			expected: "https://region.console.aws.amazon.com/codesuite/codecommit/repositories/foobar/browse/refs/heads/branch/--/test/readme.md?region=region",
		},
	}

	for _, test := range cases {
		dir := newRepo(test.remote, test.branch)
		defer os.RemoveAll(dir)

		repo, err := New(dir)
		assert.Nil(t, err, test.name)

		url, err := repo.URL(filepath.Join(dir, test.file))
		assert.Nil(t, err, test.name)

		assert.Equal(t, test.expected, url, test.name)
	}
}
