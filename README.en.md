# gh-otui

It is pronounced as /oˈtuː.i/.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)  
(All the repositories displayed in GIFs are public ones belonging to the organization I am part of.)

gh-otui is a CLI tool that combines gh, ghq, and fuzzy finders (peco, fzf).  
It allows you to search and browse the repositories of an organization using a fuzzy finder mechanism, and to clone them with ghq. This is especially useful when developing across multiple repositories, as you can complete cloning using the CLI provided you know the repository name.

## Features

- List organization repositories on GitHub
- Interactive repository selection using a fuzzy finder
- Clone selected repositories with ghq (if not already cloned)
- Visual display of cloned repositories (✓ mark)

## Prerequisite Tools

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco) 
  - or [fzf](https://github.com/junegunn/fzf). You can use fzf by setting the environment variable `GH_OTUI_SELECTOR` to `fzf`. If the environment variable is not specified, the installed tool between peco and fzf will be used. If both are installed, peco will take priority.

## Installation

```bash
gh extension install n3xem/gh-otui
```

## How to Use

1. Create a cache of the repositories in your organization:

```bash
gh otui --cache
```

The cache will be saved in `~/.config/gh/extensions/gh-otui/cache.json`.

2. Execute the following command:

```bash
gh otui
```

3. Select the desired repository in the fuzzy finder interface
   - The ✓ mark indicates repositories that have already been cloned.
   - Selecting an un-cloned repository will trigger cloning with ghq.
   - The determination of whether a repository is cloned is done by checking the path of `ghq root`.

4. The local path of the selected repository will be outputted.
   - It is convenient to use it together with the cd command to move immediately.
   - Example: `cd $(gh otui)`

## Output Format

Repositories are displayed in the following format:

- ✓: Mark indicating a cloned repository
- organization-name: Name of the GitHub organization
- repository-name: Name of the repository