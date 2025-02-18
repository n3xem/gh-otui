# gh-otui



/oˈtuː.i/ 读作。

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
（GIF中显示的所有仓库都是我所属组织的公共资源）

gh-otui 是一个结合了 gh 和 ghq、模糊查找器（peco, fzf）的 CLI 工具。  
可以利用模糊查找器的机制，跨越组织中的仓库进行搜索和浏览，并使用 ghq 进行克隆。特别是在跨多个仓库进行开发时，只要知道仓库名称，就可以仅凭 CLI 完成克隆，非常方便。

## 功能

- 显示 GitHub 组织仓库列表
- 使用模糊查找器进行交互式的仓库选择
- 使用 ghq 克隆所选仓库（如果尚未克隆）
- 以可视化方式显示已克隆的仓库（✓ 标记）

## 前提工具

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - 或者 [fzf](https://github.com/junegunn/fzf)。通过设置环境变量 `GH_OTUI_SELECTOR` 为 `fzf`，可以使用 fzf。如果没有指定环境变量，将使用已安装的 peco 和 fzf。若两者均已安装，将优先使用 peco。
  
## 安装

```bash
gh extension install n3xem/gh-otui
```

## 使用方法

1. 创建所属组织的仓库缓存：

```bash
gh otui --cache
```

缓存将保存在 `~/.config/gh/extensions/gh-otui/cache.json`。

2. 执行以下命令：

```bash
gh otui
```

3. 在模糊查找器界面中选择目标仓库
   - ✓ 标记表示已克隆的仓库
   - 选择未克隆的仓库时，将执行 ghq 的克隆操作
   - 克隆的判断是通过检查 `ghq root` 的路径进行的

4. 所选仓库的本地路径将输出到标准输出。
   - 与 cd 命令结合使用时，可以迅速移动，非常方便。
   - 例: `cd $(gh otui)`

## 输出格式

仓库将以以下格式显示：

- ✓: 表示已克隆的仓库标记
- organization-name: GitHub 的组织名
- repository-name: 仓库名