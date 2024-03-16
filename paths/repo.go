package paths

import (
	"github.com/1f349/gomvn/database"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetRepositories(basePath string, repository []string) (map[string][]*database.Artifact, error) {
	result := map[string][]*database.Artifact{}
	for _, repo := range repository {
		result[repo] = []*database.Artifact{}
		repoPath := filepath.Join(basePath, repo)
		err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".pom") {
				path = strings.Replace(path, "\\", "/", -1)
				path = strings.Replace(path, repoPath+"/", "", 1)
				artifact := newArtifact(path, info.ModTime())
				result[repo] = append(result[repo], artifact)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func newArtifact(p string, mod time.Time) *database.Artifact {
	parts := strings.Split(p, "/")
	last := len(parts) - 1
	return &database.Artifact{
		MvnGroup: strings.Join(parts[0:last-2], "."),
		Artifact: parts[last-2],
		Version:  parts[last-1],
		Modified: mod,
	}
}
