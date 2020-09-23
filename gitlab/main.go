package gitlab

import (
	"errors"
	"fmt"
	"strings"
)

func File(remote Remote, branch, file string) (string, error) {
	if !remote.isGitlab() {
		return "", errors.New("cannot handle this remote")
	}

	remoteUrl := remote.URL()

	branch = strings.Replace(branch, "refs/heads/origin/", "", -1)
	return fmt.Sprintf("%s/-/blob/%s/%s", remoteUrl, branch, file), nil
}
