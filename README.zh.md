# gh-otui

/oˈtuː.i/ 发音。

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
（所显示的GIF库都是我所属组织的公共库）

gh-otui 是一个将 gh 和 ghq、模糊查找器（peco, fzf）结合在一起的 CLI 工具。  
可以利用模糊查找器的机制横向搜索和浏览组织或自己的库，并使用 ghq 进行克隆。特别是在横跨多个库进行开发时，只要知道库名，就能通过 CLI 完成克隆，十分方便。

## 功能

- 显示 GitHub 的组织和自己的库列表
- 使用模糊查找器进行交互式的库选择
- 对选定的库进行 ghq 克隆（如果尚未克隆）
- 以视觉方式显示已克隆的库（✓ 标记）

## 前提工具

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - 或者 [fzf](https://github.com/junegunn/fzf)。通过将环境变量 `GH_OTUI_SELECTOR` 设置为 `fzf` 以使用 fzf。如果没有指定环境变量，则使用安装了 peco 和 fzf 的工具。如果两者都已安装，优先使用 peco。
  
## 安装

```bash
gh extension install n3xem/gh-otui
```

## 使用方法

1. 只需执行 gh otui 命令。初次运行时，将创建一个缓存以存储获取的库列表。

```bash
gh otui
```

2. 在模糊查找器界面中选择目标库
   - ✓ 标记表示已克隆的库
   - 选择未克隆的库时将进行 ghq 克隆
   - 克隆状态的判断将检查 `ghq root` 的路径

3. 选定库的本地路径将输出到标准输出。
   - 与 cd 命令结合使用时，可以立即移动，十分方便。
   - 例子: `cd $(gh otui)`

## 输出格式

库将以以下格式显示：

- ✓: 表示已克隆的库的标记
- organization-name: GitHub 的组织名称
- repository-name: 库名

## 关于缓存

gh-otui 使用如下缓存结构：

- **缓存保存位置**: `~/.config/gh/extensions/gh-otui/`
- **有效期**: 1 小时（从最后更新算起1小时后被判定为过期，并在后台进行自动更新）
- **元数据文件**: `_md.json` - 保存缓存的最后更新时间
- **主机目录**: 按每个 GitHub 主机（例如：`github.com`）创建目录
- **组织文件**: 每个组织创建 `{organization}.json` 文件。保存库信息

缓存会在以下情况下进行更新：
1. 初次执行时（如果缓存不存在）
2. 缓存有效期（1小时）到期时（在后台自动更新）

缓存删除：可以使用 `gh otui clear` 命令删除缓存目录。