# gh-otui

/oˈtuː.i/ 读作。

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
（GIF中展示的所有仓库都是我所在组织的公共仓库）

gh-otui 是一个结合了 gh 和 ghq、模糊查找器（peco, fzf）的 CLI 工具。  
可以利用模糊查找器的机制横向搜索和浏览组织的仓库，并使用 ghq 进行克隆。特别是在跨多个仓库进行开发的情况下，只要知道仓库名称，就可以仅通过 CLI 完成克隆，非常方便。

## 功能

- 显示 GitHub 组织的仓库列表
- 使用模糊查找器进行交互式的仓库选择
- 选择的仓库可通过 ghq 进行克隆（如果尚未克隆的话）
- 已克隆仓库的可视化显示（✓ 标记）

## 前提工具

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - 或者 [fzf](https://github.com/junegunn/fzf)。通过将环境变量 `GH_OTUI_SELECTOR` 设置为 `fzf` 可以使用 fzf。如果没有指定环境变量，系统将使用已安装的 peco 或 fzf。如果两个都已安装，将优先使用 peco。
  
## 安装

```bash
gh extension install n3xem/gh-otui
```

## 用法

1. 创建所属组织的仓库缓存：

```bash
gh otui --cache
```

缓存将保存在 `~/.config/gh/extensions/gh-otui/cache.json` 中。

2. 运行以下命令：

```bash
gh otui
```

3. 在模糊查找器界面选择目标仓库
   - ✓ 标记表示已克隆的仓库
   - 选择未克隆的仓库时，将执行 ghq 的克隆操作
   - 克隆的判定将通过检查 `ghq root` 的路径来进行

4. 选择的仓库的本地路径将标准输出。
   - 与 cd 命令结合使用，能快速移动，十分方便。
   - 例：`cd $(gh otui)`

## 输出格式

仓库将以以下格式显示：

- ✓: 表示已克隆的仓库标记
- organization-name: GitHub 组织名
- repository-name: 仓库名