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
		exp    string
	}{
		{
			name:   "simple file",
			remote: "https://user:token@gitlab.com/team/repo",
			branch: "master",
			file:   "test/readme.md",
			exp:    "https://gitlab.com/team/repo/-/blob/master/test/readme.md",
		},
	}

	for _, test := range cases {
		r := Remote{R: test.remote}

		url, err := r.File(test.branch, test.file)
		assert.Nil(t, err, test.name)
		assert.Equal(t, test.exp, url, test.name)
	}
}
