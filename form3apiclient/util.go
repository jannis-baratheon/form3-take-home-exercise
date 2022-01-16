package form3apiclient

import (
	"net/url"
	"path"
)

func join(baseAbsoluteURL string, relativePath string) (string, error) {
	baseURL, err := url.Parse(baseAbsoluteURL)
	if err != nil {
		return "", WrapError(err, "error parsing url")
	}

	if !baseURL.IsAbs() {
		return "", URLError("baseAbsoluteURL must be absolute")
	}

	if len(baseURL.Query()) > 0 {
		return "", URLError("baseAbsoluteURL with query is not supported")
	}

	if baseURL.Fragment != "" {
		return "", URLError("baseAbsoluteURL with fragment is not supported")
	}

	relativeURL, err := url.Parse(relativePath)
	if err != nil {
		return "", WrapError(err, "error parsing url")
	}

	if len(relativeURL.Query()) > 0 {
		return "", URLError("relativeURL with query is not supported")
	}

	if relativeURL.Fragment != "" {
		return "", URLError("relativeURL with fragment is not supported")
	}

	baseURL.Path = path.Join(baseURL.Path, relativePath)

	return baseURL.String(), nil
}
