#compdef todoist

_cli_zsh_autocomplete() {

  _arguments -C \
    '1:cmd:->cmds' \
    '2:close:->tasks' \
    '*:: :->args' \
  && ret=0

  case "$state" in
    (cmds)
      local -a commands;
      commands=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
      _describe -t commands 'command' commands && ret=0
    ;;
    (tasks)
      local -a tasks; tasks=("${(@f)$(todoist l | sed 's/\s/:/')}")
      _describe -t tasks 'task' tasks && ret=0
    ;;
  esac

  return 1
}

compdef _cli_zsh_autocomplete todoist
