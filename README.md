# gh-otui

gh-otuiは、GitHubのOrganizationのリポジトリを簡単に閲覧・クローンし、クローン先のpathを出力するCLIツールです。
特に複数のリポジトリを横断して開発している場合、リポジトリ名さえ知っていればCLIのみでクローンを完結できるので便利です.
使う際にキャッシュを行うことで、初回のみGitHub APIのリクエストが発生するので、2回目以降は高速に動作します.

gh-otui = gh + org + tui

## 機能

- GitHubの組織リポジトリの一覧表示
- pecoを使用した対話的なリポジトリ選択
- 選択したリポジトリの自動クローン（未クローンの場合）
- クローン済みリポジトリの視覚的な表示（✓マーク）

## 必要条件

`make deps` を実行することでまとめてbrewでインストールできます。

- [GitHub CLI](https://cli.github.com/) (gh)
- [peco](https://github.com/peco/peco)
- [ghq](https://github.com/x-motemen/ghq)

## 使い方

1. GitHub CLIにログインしていることを確認してください：

```bash
gh auth login
```

2. 所属しているorganizationのリポジトリのキャッシュを作成します：

```bash
go run ./main.go --cache
```

3. 以下のコマンドでリポジトリを選択します：

```bash
go run ./main.go
```

4. pecoインターフェースで目的のリポジトリを選択します
   - ✓マークは既にクローン済みのリポジトリを示します
   - 未クローンのリポジトリを選択すると自動的にクローンされます

5. 選択したリポジトリのローカルパスが出力されます。
   - cdコマンドと連携して使用するとすぐ移動できて便利です。


## 出力形式

リポジトリは以下の形式で表示されます：
- ✓: クローン済みリポジトリを示すマーク
- organization-name: GitHubの組織名
- repository-name: リポジトリ名
