package gitlab

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

func File(remote Remote, dir, branch, file string) (string, error) {
	if !remote.isGitlab() {
		return "", errors.New("cannot handle this remote")
	}

	absFile, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}

	relativeFile := strings.TrimPrefix(strings.Replace(absFile, dir, "", -1), "/")
	fmt.Println("relative file", relativeFile)
	remoteUrl := remote.URL()

	branch = strings.Replace(branch, "refs/heads/origin/", "", -1)
	return fmt.Sprintf("%s/-/blob/%s/%s", remoteUrl, branch, relativeFile), nil
}
