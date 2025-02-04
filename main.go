package main

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
)

func main() {
	client, err := api.DefaultRESTClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 組織一覧を格納する構造体のスライスを定義
	var orgs []struct {
		Login string `json:"login"`
	}

	// GET /user/orgs エンドポイントを呼び出して組織一覧を取得
	err = client.Get("user/orgs", &orgs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("所属している組織とリポジトリ一覧:")
	for _, org := range orgs {
		fmt.Printf("Organization: %s\n", org.Login)

		// 組織のリポジトリ一覧を格納する構造体のスライスを定義
		var repos []struct {
			Name string `json:"name"`
		}

		// GET /orgs/{org}/repos エンドポイントを呼び出してリポジトリ一覧を取得
		err = client.Get(fmt.Sprintf("orgs/%s/repos", org.Login), &repos)
		if err != nil {
			fmt.Printf("  リポジトリの取得に失敗: %v\n", err)
			continue
		}

		// リポジトリ一覧を表示
		for _, repo := range repos {
			fmt.Printf("  - %s\n", repo.Name)
		}
		fmt.Println()
	}
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
