// Package git 数据提取层
package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type BranchInfo struct {
	Name   string
	Commit string
	Parent []string
}

func LoadBranches(repoPath string) ([]BranchInfo, error) {
	// 打开本地仓库
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	var result []BranchInfo

	// 遍历所有分支引用
	refs, err := r.References()
	if err != nil {
		return nil, err
	}

	//在 go-git 中，plumbing.Reference 表示 Git 中的一个「引用」，比如：
	//1.refs/heads/main（本地分支）
	//2.refs/remotes/origin/dev（远程分支）
	//3.refs/tags/v1.0.0（标签）
	//4.HEAD（当前 HEAD）

	// plumbing.Reference 结构体
	//	type Reference struct {
	//		name ReferenceName // 引用名，比如 refs/heads/main
	//		hash Hash          // 目标对象的哈希（通常是 commit）
	//	}
	refs.ForEach(func(ref *plumbing.Reference) error {
		// 这里的 ref 就是一个指向某个对象（通常是 commit）的引用
		fmt.Println("ref------", ref)
		// 判断是否是分支引用
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
