package gitlab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	var cases = []struct {
		name         string
		remote       Remote
		branch       string
		relativeFile string
		exp          string
	}{
		{
			name:         "file inside the repo",
			remote:       Remote{"https://user:token@gitlab.com/foo/bar"},
			branch:       "master",
			relativeFile: "folder/readme.md",
			exp:          "https://gitlab.com/foo/bar/-/blob/master/folder/readme.md",
		},
	}

	for _, test := range cases {
		remote, err := File(test.remote, test.branch, test.relativeFile)
		assert.Nil(t, err)
		assert.Equal(t, test.exp, remote, test.name)
	}
}
