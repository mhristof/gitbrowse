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
