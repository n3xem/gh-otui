# gh-otui

<!-- ss-markdown-ignore start -->
[English](README.en.md) | [简体中文](README.zh.md) | [Español](README.es.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [한국어](README.ko.md)
<!-- ss-markdown-ignore end -->

/oˈtuː.i/ と読みます。

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
（GIFで表示されているリポジトリはすべて自分が所属している組織のパブリックのものです）

gh-otuiはghとghq、fuzzy finder (peco, fzf)を組み合わせたCLIツールです。
Organizationや自分のリポジトリをfuzzy finderの仕組みを使って横断して検索・閲覧し、ghqでクローンすることができます。特に複数のリポジトリを横断して開発している場合、リポジトリ名さえ知っていればCLIのみでクローンを完結できるので便利です.

## 機能

- GitHubのOrganization、自分のリポジトリの一覧表示
- fuzzy finderを使用した対話的なリポジトリ選択
- 選択したリポジトリのghqによるクローン（未クローンの場合）
- クローン済みリポジトリの視覚的な表示（✓マーク）

## 前提ツール

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - または [fzf](https://github.com/junegunn/fzf)。環境変数 `GH_OTUI_SELECTOR` を `fzf` に設定することでfzfを使用できます。環境変数の指定がない場合は、pecoとfzfのインストールされている方を使います。両方インストールされている場合はpecoが優先されます。

## インストール

```bash
gh extension install qawatake/gh-otui
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

3. fuzzy finderインターフェースで目的のリポジトリを選択します
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
