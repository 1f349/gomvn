package paths

import (
	"github.com/1f349/gomvn/database/types"
	"os"
	"path/filepath"
	"strings"
)

func GetRepositories(basePath string, repository []string) map[string][]*types.Artifact {
	result := map[string][]*types.Artifact{}
	for _, repo := range repository {
		result[repo] = []*types.Artifact{}
		repoPath := filepath.Join(basePath, repo)
		_ = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".pom") {
				path = strings.Replace(path, "\\", "/", -1)
				path = strings.Replace(path, repoPath+"/", "", 1)

				parts := strings.Split(path, "/")
				last := len(parts) - 1
				artifact := &types.Artifact{
					MvnGroup: strings.Join(parts[0:last-2], "."),
					Artifact: parts[last-2],
					Version:  parts[last-1],
					Modified: info.ModTime(),
				}
				result[repo] = append(result[repo], artifact)
			}
			return nil
		})
	}
	return result
}
