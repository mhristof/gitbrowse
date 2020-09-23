package github

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
			name:   "valid https remote with .git",
			remote: "https://github.com/user/repo.git",
			exp:    "https://github.com/user/repo",
		},
		{
			name:   "valid https remote without .git",
			remote: "https://github.com/user/repo",
			exp:    "https://github.com/user/repo",
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
			remote: "https://github.com/user/repo",
			branch: "master",
			file:   "test/readme.md",
			exp:    "https://github.com/user/repo/blob/master/test/readme.md",
		},
		{
			name:   "simple file",
			remote: "https://github.com/user/repo",
			branch: "foobar",
			file:   "test/readme.md",
			exp:    "https://github.com/user/repo/blob/foobar/test/readme.md",
		},
	}

	for _, test := range cases {
		r := Remote{R: test.remote}

		url, err := r.File(test.branch, test.file)
		assert.Nil(t, err, test.name)
		assert.Equal(t, test.exp, url, test.name)
	}
}
