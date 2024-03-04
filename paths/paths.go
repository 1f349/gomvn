package paths

import (
	"fmt"
	"net/http"
	"strings"
)

type Paths struct {
	Repository []string
}

func (p Paths) NormalizePath(path string) string {
	if path[0] == '/' {
		path = path[1:]
	}
	if strings.Contains(path, "..") || strings.Contains(path, "~") {
		return ""
	}
	if strings.Count(path, "/") <= 1 {
		return path
	}
	for _, repo := range p.Repository {
		if strings.HasPrefix(path, repo) {
			return path
		}
	}
	return ""
}

func (p Paths) ParsePath(req *http.Request) (string, error) {
	path := p.NormalizePath(req.URL.Path)
	if strings.Count(path, "/") < 3 {
		return "", fmt.Errorf("path should be repository/group/artifact")
	}
	return path, nil
}

func (p Paths) ParsePathParts(req *http.Request) (string, string, string, error) {
	path, err := p.ParsePath(req)
	if err != nil {
		return "", "", "", err
	}
	parts := strings.Split(path, "/")
	last := len(parts) - 1
	return parts[0], strings.Join(parts[1:last-1], "/"), parts[last], nil
}
