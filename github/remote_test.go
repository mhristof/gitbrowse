package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	cases := []struct {
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
		{
			name:   "valid git remote",
			remote: "git@github.com:user/repo.git",
			exp:    "https://github.com/user/repo",
		},
	}

	for _, test := range cases {
		r := Remote{test.remote}
		assert.Equal(t, test.exp, r.URL(), test.name)
	}
}

func TestFile(t *testing.T) {
	cases := []struct {
		name   string
		remote string
		branch string
		file   string
		line   int
		exp    string
	}{
		{
			name:   "simple file",
			remote: "https://github.com/user/repo",
			branch: "master",
			file:   "test/readme.md",
			line:   -1,
			exp:    "https://github.com/user/repo/tree/master/test/readme.md",
		},
		{
			name:   "simple file",
			remote: "https://github.com/user/repo",
			branch: "foobar",
			file:   "test/readme.md",
			line:   -1,
			exp:    "https://github.com/user/repo/tree/foobar/test/readme.md",
		},
		{
			name:   "simple file with line number",
			remote: "https://github.com/user/repo",
			branch: "foobar",
			file:   "test/readme.md",
			line:   100,
			exp:    "https://github.com/user/repo/tree/foobar/test/readme.md#L100",
		},
	}

	for _, test := range cases {
		r := Remote{R: test.remote}

		url, err := r.File(test.branch, test.file, test.line)
		assert.Nil(t, err, test.name)
		assert.Equal(t, test.exp, url, test.name)
	}
}

func TestGitToHTTP(t *testing.T) {
	cases := []struct {
		name   string
		remote string
		exp    string
	}{
		{
			name:   "git remote",
			remote: "git@github.com:mhristof/alfred-pbpaste.git",
			exp:    "https://github.com/mhristof/alfred-pbpaste.git",
		},
		{
			name:   "http remote (passthrough mode)",
			remote: "https://github.com/mhristof/alfred-pbpaste.git",
			exp:    "https://github.com/mhristof/alfred-pbpaste.git",
		},
	}

	for _, test := range cases {
		assert.Equal(t, test.exp, gitToHTTP(test.remote), test.name)
	}
}
