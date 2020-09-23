package git

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newRepo(remote string) string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("cd %s && git init", dir))
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("git", "-C", dir, "remote", "add", "origin", remote)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	return dir
}

func TestURL(t *testing.T) {
	var cases = []struct {
		name     string
		remote   string
		file     string
		path     string
		expected string
	}{
		{
			name:     "gitlab https remote with username and token",
			remote:   "https://glUsernames:glToken@gitlab.com/project/repository",
			file:     "README.md",
			expected: "https://gitlab.com/project/repository/-/blob/refs/heads/master/README.md",
		},
	}

	for _, test := range cases {
		dir := newRepo(test.remote)

		repo, err := New(dir)
		fmt.Println(err)
		assert.Nil(t, err, test.name)

		url, err := repo.URL(filepath.Join(dir, test.file))
		assert.Nil(t, err, test.name)

		assert.Equal(t, test.expected, url, test.name)
	}
}
