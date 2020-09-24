package gitlab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	var cases = []struct {
		name   string
		remote string
		exp    string
	}{
		{
			name:   "valid https remote",
			remote: "https://username:token@gitlab.com/foo/bar",
			exp:    "https://gitlab.com/foo/bar",
		},
	}

	for _, test := range cases {
		r := Remote{test.remote}
		assert.Equal(t, test.exp, r.URL(), test.name)
	}
}

func TestFile(t *testing.T) {
	var cases = []struct {
		name   string
		remote string
		branch string
		file   string
		line   int
		exp    string
	}{
		{
			name:   "simple file",
			remote: "https://user:token@gitlab.com/team/repo",
			branch: "master",
			file:   "test/readme.md",
			line:   -1,
			exp:    "https://gitlab.com/team/repo/-/blob/master/test/readme.md",
		},
		{
			name:   "simple file with line",
			remote: "https://user:token@gitlab.com/team/repo",
			branch: "master",
			file:   "test/readme.md",
			line:   100,
			exp:    "https://gitlab.com/team/repo/-/blob/master/test/readme.md#L100",
		},
	}

	for _, test := range cases {
		r := Remote{R: test.remote}

		url, err := r.File(test.branch, test.file, test.line)
		assert.Nil(t, err, test.name)
		assert.Equal(t, test.exp, url, test.name)
	}
}
