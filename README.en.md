# gh-otui

It is pronounced as /oˈtuː.i/.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(All repositories displayed in GIF are public ones belonging to the organization I am part of.)

gh-otui is a CLI tool that combines gh and ghq with a fuzzy finder (peco, fzf).  
You can search and browse across organizations and your own repositories using the fuzzy finder mechanism and clone them with ghq. This is particularly convenient when developing across multiple repositories, as you can complete cloning using only the CLI if you know the repository name.

## Features

- Lists GitHub organizations and your own repositories
- Interactive repository selection using a fuzzy finder
- Clones the selected repository using ghq (if it is not already cloned)
- Visual display of cloned repositories (✓ mark)

## Prerequisite Tools

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - Or [fzf](https://github.com/junegunn/fzf). You can use fzf by setting the environment variable `GH_OTUI_SELECTOR` to `fzf`. If the environment variable is not set, it will use whichever of peco or fzf is installed. If both are installed, peco takes precedence.

## Installation

```bash
gh extension install n3xem/gh-otui
```

## Usage

1. Create a cache of the repositories in your organization:

```bash
gh otui --cache
```

The cache will be saved to `~/.config/gh/extensions/gh-otui/cache.json`.

2. Run the following command:

```bash
gh otui
```

3. Select the desired repository in the fuzzy finder interface:
   - The ✓ mark indicates repositories that have already been cloned.
   - Selecting an un-cloned repository will initiate cloning via ghq.
   - The determination of whether a repository has been cloned is based on checking the path of `ghq root`.

4. The local path of the selected repository will be output to the standard output.
   - It is convenient to use this in conjunction with the cd command for quick navigation.
   - Example: `cd $(gh otui)`

## Output Format

Repositories are displayed in the following format:

- ✓: Mark indicating a cloned repository
- organization-name: Name of the GitHub organization
- repository-name: Name of the repository