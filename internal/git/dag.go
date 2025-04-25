package git

// internal/git/dag.go
import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type CommitNode struct {
	Hash     string   `json:"hash"`
	Parents  []string `json:"parents"`
	Children []string `json:"children"`
	Branches []string `json:"branches"` // 所属分支名
}

func BuildCommitDAG(repoPath string) (map[string]*CommitNode, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	commitGraph := make(map[string]*CommitNode)

	// 遍历所有分支
	refs, err := r.References()
	if err != nil {
		return nil, err
	}

	refs.ForEach(func(ref *plumbing.Reference) error {
		if !ref.Name().IsBranch() {
			return nil
		}
		branchName := ref.Name().Short()

		iter, err := r.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			return err
		}

		// 遍历分支上的每个 commit
		err = iter.ForEach(func(c *object.Commit) error {
			node, exists := commitGraph[c.Hash.String()]
			if !exists {
				node = &CommitNode{
					Hash:     c.Hash.String(),
					Parents:  []string{},
					Children: []string{},
					Branches: []string{},
				}
				commitGraph[c.Hash.String()] = node
			}

			node.Branches = appendIfMissing(node.Branches, branchName)

			for _, p := range c.ParentHashes {
				parentHash := p.String()

				// 添加到当前 node 的 parents
				node.Parents = appendIfMissing(node.Parents, parentHash)

				// 添加到 parent 的 children
				parent, ok := commitGraph[parentHash]
				if !ok {
					parent = &CommitNode{
						Hash:     parentHash,
						Parents:  []string{},
						Children: []string{},
						Branches: []string{},
					}
					commitGraph[parentHash] = parent
				}
				parent.Children = appendIfMissing(parent.Children, node.Hash)
			}
			return nil
		})

		return err
	})

	return commitGraph, nil
}

func appendIfMissing(list []string, item string) []string {
	for _, v := range list {
		if v == item {
			return list
		}
	}
	return append(list, item)
}
