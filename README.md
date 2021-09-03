Todoist CLI client
===

[Todoist](https://todoist.com/) CLI Client, written in Golang.

## Description

[Todoist](https://todoist.com/) is a cool TODO list web application.
This program will let you use the Todoist in CLI.

![color image](https://cloud.githubusercontent.com/assets/6121271/20603278/2261b424-b2a4-11e6-8fa7-d533e2144942.png)

## Demo (with [peco](https://github.com/peco/peco))

### Add Task

![Add task](https://cloud.githubusercontent.com/assets/6121271/19836528/6ed99956-9ee6-11e6-85b0-7539393d803b.gif)

### Close Task

![Close task](https://cloud.githubusercontent.com/assets/6121271/19836531/7c399218-9ee6-11e6-974c-9dd59ced13a5.gif)

## Usage

```
$ todoist --help
NAME:
   todoist - Todoist CLI Client

USAGE:
   todoist [global options] command [command options] [arguments...]

VERSION:
   0.15.0

COMMANDS:
     list, l                  Show all tasks
     show                     Show task detail
     completed-list, c-l, cl  Show all completed tasks (only premium users)
     add, a                   Add task
     modify, m                Modify task
     close, c                 Close task
     delete, d                Delete task
     labels                   Show all labels
     projects                 Show all projects
     karma                    Show karma
     sync, s                  Sync cache
     quick, q                 Quick add a task
     help, h                  Show a list of commands or help for one command

GLOBAL OPTIONS:
   --color              colorize output
   --csv                output in CSV format
   --debug              output logs
   --namespace          display parent task like namespace
   --indent             display children task with indent
   --project-namespace  display parent project like namespace
   --help, -h           show help
   --version, -v        print the version
```

### `list --filter`

You can filter tasks by `--filter` option on `list` subcommand.
The filter syntax is base on [todoist official filter syntax](https://support.todoist.com/hc/en-us/articles/205248842-Filters).

Supported filter is [here](https://github.com/sachaos/todoist/issues/15#issuecomment-334140101).

#### e.g. List tasks which over due date and have high priority

```
todoist list --filter '(overdue | today) & p1'
```

## Config

Config by default stored in `$HOME/.config/todoist/config.json`

It has following parameters:

```
{
  "token": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", # todoist api token, required
  "color": "true"                                      # colorize all output, not required, default false
}
```

## Install

### Homebrew (Mac OS)

```
$ brew tap sachaos/todoist
$ brew install todoist
```

### AUR

* [todoist](https://aur.archlinux.org/packages/todoist/)
* [todoist-git](https://aur.archlinux.org/packages/todoist-git/)

### Nix/NixOS

```
nix-env -iA nixos.todoist
```

It's important to notice that if you're using NixOS, the cache and config file will be present at your home directory: `~/.todoist.cache.json` and `~/.todoist.config.json`.

### Docker

```
$ git clone https://github.com/sachaos/todoist.git
$ cd todoist
$ make docker-build token=xxxxxxxxxxxxxxxxxxxx
$ make docker-run
```

You will be running the next commands from inside the container.

PS: We add a step that is run `sync` before any command, so you will be always up to date!

### Build it yourself

You need go 1.12.

```
$ mkdir -p $GOPATH/src/github.com/sachaos
$ cd $GOPATH/src/github.com/sachaos
$ git clone https://github.com/sachaos/todoist.git
$ cd todoist
$ make install
```

### Register API token

When you run `todoist` first time, you will be asked your Todoist API token.
Please input Todoist API token and register it. In order to get your API token
go to [https://todoist.com/prefs/integrations](https://todoist.com/prefs/integrations)

### Sync

After register API token, you should sync with todoist.com by `sync` sub command, like below.

```
$ todoist sync
```

### Use with peco/fzf

**RECOMMENDED**

Install [peco](https://github.com/peco/peco) and load `todoist_functions.sh` on your `.zshrc`, like below.

fish version is here. [ka2n/fish-peco_todoist](https://github.com/ka2n/fish-peco_todoist) Thanks @ka2n!

If you would prefer to use [fzf](https://github.com/junegunn/fzf) instead load `todoist_functions_fzf.sh` like below.

```
$ source "$GOPATH/src/github.com/sachaos/todoist/todoist_functions.sh"
```

#### If installed via homebrew

If installed via homebrew and using zsh (usually this is added to your `.zshrc` without the `$`, usually before loading your ZSH plugin manager):

For **peco**:
```
$ source $(brew --prefix)/share/zsh/site-functions/_todoist_peco
```

For **fzf**:
```
$ source $(brew --prefix)/share/zsh/site-functions/_todoist_fzf
```

**TODO**: fish + homebrew

#### keybind

```
<C-x> t t: select task with peco
<C-x> t p: select project with peco
<C-x> t l: select labels with peco
<C-x> t c: select task and close with peco
<C-x> t d: select date
<C-x> t o: select task, and open it with browser when has url
```

### Enable shell completion

You can also enable shell completion by adding the following lines to your `.bashrc`/`.zshrc` files.

```
# Bash
PROG=todoist source "$GOPATH/src/github.com/urfave/cli/autocomplete/bash_autocomplete"
# Zsh
PROG=todoist source "$GOPATH/src/github.com/urfave/cli/autocomplete/zsh_autocomplete"
```

## Author

[sachaos](https://github.com/sachaos)
