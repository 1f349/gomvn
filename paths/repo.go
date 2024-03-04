package paths

import (
	"os"
	"path/filepath"
	"strings"
)

func GetRepositories(basePath string, repository []string) map[string][]*entity.Artifact {
	result := map[string][]*database.Artifact{}
	for _, repo := range repository {
		result[repo] = []*database.Artifact{}
		repoPath := filepath.Join(basePath, repo)
		_ = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".pom") {
				path = strings.Replace(path, "\\", "/", -1)
				path = strings.Replace(path, repoPath+"/", "", 1)
				artifact := entity.NewArtifact(path, info.ModTime())
				result[repo] = append(result[repo], artifact)
			}
			return nil
		})
	}
	return result
}
