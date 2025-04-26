package main

import (
	"encoding/json"
	"github.com/henry-insomniac/branch-manager/internal/git"
	"log"
	"os"
)

func main() {
	path := "/Users/mm/work-space/rc-workShop"
	dag, err := git.BuildCommitDAG(path)
	branchs, err := git.LoadBranches(path)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatalf("构建 DAG 失败: %v", err)
	}

	// 转成 JSON 打印
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(dag)
	_ = enc.Encode(branchs)
}
