# gh-otui

It is read as /oˈtuː.i/.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(All repositories displayed in the GIF are public ones belonging to the organization I am part of.)

gh-otui is a CLI tool that combines gh, ghq, and fuzzy finders (peco, fzf).  
It allows you to search and browse through your organization's repositories using a fuzzy finder mechanism, and you can clone them with ghq. This is especially useful when developing across multiple repositories, as you can complete cloning using only the CLI if you know the repository name.

## Features

- Display a list of GitHub organization repositories
- Interactive repository selection using a fuzzy finder
- Clone selected repositories using ghq (if not already cloned)
- Visual display of cloned repositories (✓ mark)

## Prerequisite Tools

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - Or [fzf](https://github.com/junegunn/fzf). You can use fzf by setting the environment variable `GH_OTUI_SELECTOR` to `fzf`. If no environment variable is specified, the installed tool between peco and fzf will be used. If both are installed, peco takes priority.
  
## Installation

```bash
gh extension install n3xem/gh-otui
```

## Usage

1. Create a cache of the repositories in your organization:

```bash
gh otui --cache
```

The cache is saved at `~/.config/gh/extensions/gh-otui/cache.json`.

2. Execute the following command:

```bash
gh otui
```

3. Select the desired repository in the fuzzy finder interface:
   - The ✓ mark indicates that the repository is already cloned.
   - If you select an un-cloned repository, it will be cloned using ghq.
   - Cloning status is determined by checking the path of `ghq root`.

4. The local path of the selected repository will be printed to standard output.
   - It’s convenient to use this in conjunction with the cd command for immediate navigation.
   - Example: `cd $(gh otui)`

## Output Format

Repositories are displayed in the following format:

- ✓: Mark indicating a cloned repository
- organization-name: Name of the GitHub organization
- repository-name: Name of the repository