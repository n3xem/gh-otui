# gh-otui



/oˈtuː.i/ 读作。

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
（GIF中显示的所有仓库都是我所属组织的公共仓库）

gh-otui是一个结合了gh和ghq、模糊查找工具（peco, fzf）的CLI工具。  
可以通过模糊查找机制横向搜索和浏览组织及自己的仓库，并使用ghq进行克隆。特别是在跨多个仓库进行开发时，只要知道仓库名，就可以仅通过CLI完成克隆，十分方便。

## 功能

- 显示GitHub的组织和自己仓库的列表
- 使用模糊查找器进行交互式仓库选择
- 对选择的仓库进行ghq克隆（如果未克隆）
- 已克隆仓库的视觉显示（✓标记）

## 前提工具

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - 或者 [fzf](https://github.com/junegunn/fzf)。通过将环境变量 `GH_OTUI_SELECTOR` 设置为 `fzf`，可以使用fzf。如果没有指定环境变量，将使用已安装的peco或fzf。如果两者都已安装，将优先使用peco。
  
## 安装

```bash
gh extension install n3xem/gh-otui
```

## 使用方法

1. 创建您所属的组织的仓库缓存：

```bash
gh otui --cache
```

缓存将保存在 `~/.config/gh/extensions/gh-otui/cache.json` 中。

2. 执行以下命令：

```bash
gh otui
```

3. 在模糊查找器界面中选择目标仓库
   - ✓标记表示已克隆的仓库
   - 选择未克隆的仓库将进行ghq克隆
   - 克隆状态的判定通过检查 `ghq root` 的路径进行

4. 选定仓库的本地路径将被输出到标准输出。
   - 与cd命令结合使用可以方便地立即移动。
   - 例： `cd $(gh otui)`

## 输出格式

仓库将以以下格式显示：

- ✓: 表示已克隆的仓库的标记
- organization-name: GitHub的组织名
- repository-name: 仓库名