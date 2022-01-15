package form3apiclient

import (
	"fmt"
	"net/url"
	"path"
)

func join(baseAbsoluteURL string, relativePath string) (string, error) {
	baseURL, err := url.Parse(baseAbsoluteURL)
	if err != nil {
		return "", err
	}

	if !baseURL.IsAbs() {
		return "", fmt.Errorf("baseAbsoluteURL must be absolute")
	}

	if len(baseURL.Query()) > 0 {
		return "", fmt.Errorf("baseAbsoluteURL with query is not supported")
	}

	if baseURL.Fragment != "" {
		return "", fmt.Errorf("baseAbsoluteURL with fragment is not supported")
	}

	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		return "", err
	}

	if len(relativeURL.Query()) > 0 {
		return "", fmt.Errorf("relativeURL with query is not supported")
	}

	if relativeURL.Fragment != "" {
		return "", fmt.Errorf("relativeURL with fragment is not supported")
	}

	baseURL.Path = path.Join(baseURL.Path, relativePath)

	return baseURL.String(), nil
}
