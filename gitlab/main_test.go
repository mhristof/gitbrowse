package gitlab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	var cases = []struct {
		name   string
		remote Remote
		dir    string
		branch string
		file   string
		exp    string
	}{
		{
			name:   "file inside the repo",
			remote: Remote{"https://user:token@gitlab.com/foo/bar"},
			dir:    "/code",
			branch: "master",
			file:   "/code/readme",
			exp:    "https://gitlab.com/foo/bar/-/blob/master/readme",
		},
	}

	for _, test := range cases {
		remote, err := File(test.remote, test.dir, test.branch, test.file)
		assert.Nil(t, err)
		assert.Equal(t, test.exp, remote, test.name)
	}
}
