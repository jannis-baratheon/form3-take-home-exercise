package form3apiclient

import (
	"fmt"
	"net/url"
	"path"
)

func join(baseAbsoluteUrl string, relativePath string) (string, error) {
	url, err := url.Parse(baseAbsoluteUrl)

	if err != nil {
		return "", err
	}

	if !url.IsAbs() {
		return "", fmt.Errorf("api url must be absolute")
	}

	url.Path = path.Join(url.Path, relativePath)
	return url.String(), nil
}