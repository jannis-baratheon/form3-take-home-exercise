package form3apiclient

import (
	"fmt"
	"net/url"
	"path"
)

func join(baseAbsoluteUrl string, relativePath string) (string, error) {
	baseUrl, err := url.Parse(baseAbsoluteUrl)

	if err != nil {
		return "", err
	}

	if !baseUrl.IsAbs() {
		return "", fmt.Errorf("baseAbsoluteUrl must be absolute")
	}

	if len(baseUrl.Query()) > 0 {
		return "", fmt.Errorf("baseAbsoluteUrl with query is not supported")
	}

	if baseUrl.Fragment != "" {
		return "", fmt.Errorf("baseAbsoluteUrl with fragment is not supported")
	}

	relativeUrl, err := url.Parse(relativePath)

	if err != nil {
		return "", err
	}

	if len(relativeUrl.Query()) > 0 {
		return "", fmt.Errorf("relativeUrl with query is not supported")
	}

	if relativeUrl.Fragment != "" {
		return "", fmt.Errorf("relativeUrl with fragment is not supported")
	}

	baseUrl.Path = path.Join(baseUrl.Path, relativePath)
	return baseUrl.String(), nil
}