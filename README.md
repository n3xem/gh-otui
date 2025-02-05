# gh-otui

/oˈtuː.i/ と読みます。

gh-otui = gh + org + tui



![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
（GIFで表示されているリポジトリはすべて自分が所属している組織のパブリックのものです）




gh-otuiはghとghq、pecoを組み合わせたCLIツールです。  
Organizationのリポジトリをpecoの仕組みを使って横断して検索・閲覧し、ghqでクローンすることができます。特に複数のリポジトリを横断して開発している場合、リポジトリ名さえ知っていればCLIのみでクローンを完結できるので便利です.  

## 機能

- GitHubの組織リポジトリの一覧表示
- pecoを使用した対話的なリポジトリ選択
- 選択したリポジトリのghqによるクローン（未クローンの場合）
- クローン済みリポジトリの視覚的な表示（✓マーク）

## 必要条件

`make deps` を実行することでまとめてbrewでインストールできます。

- [GitHub CLI](https://cli.github.com/) (gh)
- [peco](https://github.com/peco/peco)
- [ghq](https://github.com/x-motemen/ghq)


## インストール
```bash
gh extension install n3xem/gh-otui
```

## 使い方
1. 所属しているorganizationのリポジトリのキャッシュを作成します：

```bash
gh otui --cache
```

キャッシュは `~/.config/gh/extensions/gh-otui/cache.json` に保存されます。

2. 以下のコマンドを実行します：

```bash
gh otui
```

3. pecoインターフェースで目的のリポジトリを選択します
   - ✓マークは既にクローン済みのリポジトリを示します
   - 未クローンのリポジトリを選択するとghqによるクローンが行われます
   - クローン済みの判定は `ghq root` のパスを確認して行われます

4. 選択したリポジトリのローカルパスが標準出力されます。
   - cdコマンドと連携して使用するとすぐ移動できて便利です。
   - 例: `cd $(gh otui)`


## 出力形式

リポジトリは以下の形式で表示されます：
- ✓: クローン済みリポジトリを示すマーク
- organization-name: GitHubの組織名
- repository-name: リポジトリ名
