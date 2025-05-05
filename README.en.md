# gh-otui

It is read as /oˈtuː.i/.

gh-otui = gh + org + tui

![otui](https://github.com/user-attachments/assets/0c7626eb-c639-4f4c-86e1-b4ba6dab5bec)
(All repositories displayed in the GIF belong to public ones within the organization I belong to.)

gh-otui is a CLI tool that combines gh, ghq, and fuzzy finders (peco, fzf).  
You can traverse and search/view organizations and your own repositories using the fuzzy finder mechanism, and clone them with ghq. It is particularly convenient when developing across multiple repositories because if you know the repository name, you can complete the cloning using only the CLI.

## Features

- List of GitHub organizations and your own repositories
- Interactive repository selection using a fuzzy finder
- Cloning of the selected repository with ghq (if it is not already cloned)
- Visual indication of already cloned repositories (✓ mark)

## Prerequisite Tools

- [GitHub CLI](https://cli.github.com/) (gh)
- [ghq](https://github.com/x-motemen/ghq)
- [peco](https://github.com/peco/peco)
  - Or [fzf](https://github.com/junegunn/fzf). By setting the environment variable `GH_OTUI_SELECTOR` to `fzf`, you can use fzf. If no environment variable is specified, it will use whichever is installed, peco or fzf. If both are installed, peco takes precedence.
  
## Installation

```bash
gh extension install n3xem/gh-otui
```

## Usage

1. Just run the `gh otui` command. A cache containing the list of repositories to be fetched will be created for the first time.

```bash
gh otui
```

2. Select the desired repository from the fuzzy finder interface.
   - The ✓ mark indicates a repository that has already been cloned.
   - Selecting an un-cloned repository will result in a clone via ghq.
   - Cloning status is determined by checking the path of `ghq root`.

3. The local path of the selected repository will be printed to standard output.
   - It is convenient when used in conjunction with the `cd` command for quick navigation.
   - Example: `cd $(gh otui)`

## Output Format

Repositories will be displayed in the following format:

- ✓: A mark indicating a cloned repository
- organization-name: GitHub organization name
- repository-name: Repository name

## About Cache

gh-otui uses the following cache structure:

- **Cache Storage Location**: `~/.config/gh/extensions/gh-otui/`
- **Validity Period**: 1 hour (after 1 hour since the last update, it is deemed old and will be automatically updated in the background)
- **Metadata File**: `_md.json` - saves the last update time of the cache
- **Host Directory**: Creates a directory for each GitHub host (e.g., `github.com`)
- **Organization Files**: Creates a `{organization}.json` file for each organization to store repository information

The cache will be updated in the following cases:
1. On first execution (if the cache does not exist)
2. When the cache validity period (1 hour) expires (will automatically update in the background)

To delete the cache: You can delete the cache directory using the `gh otui clear` command.