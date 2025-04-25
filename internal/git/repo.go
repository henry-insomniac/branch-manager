// internal/git/repo.go
package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type BranchInfo struct {
	Name   string
	Commit string
	Parent []string
}

func LoadBranches(repoPath string) ([]BranchInfo, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	var result []BranchInfo

	refs, err := r.References()
	if err != nil {
		return nil, err
	}

	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			commit, err := r.CommitObject(ref.Hash())
			if err != nil {
				return err
			}

			var parents []string
			for _, p := range commit.ParentHashes {
				parents = append(parents, p.String())
			}

			result = append(result, BranchInfo{
				Name:   ref.Name().Short(),
				Commit: commit.Hash.String(),
				Parent: parents,
			})
		}
		return nil
	})

	return result, nil
}
