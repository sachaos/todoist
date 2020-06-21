#!/usr/bin/env bash
# Bash completion for https://github.com/sachaos/todoist

_todoist() {
    local i cur prev opts cmd fzfquery
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts=''
    cmd=''
    fzfquery=
    [ -n "$cur" ] && fzfquery="-q $cur"

    for i in "${COMP_WORDS[@]}"; do
        case "${i}" in
        todoist)
            cmd='todoist'
            ;;
        # These are the current commands; not all have completion options,
        # but they're listed here anyway, for the future
        list|show|completed-list|add|modify|close|delete|labels|projects|\
        karma|sync|quick|help)
            cmd+="__${i}"
            ;;
        l)
            cmd+='__list'
            ;;
        c-l|cl)
            cmd+='__completed-list'
            ;;
        a)
            cmd+='__add'
            ;;
        m)
            cmd+='__modify'
            ;;
        c)
            cmd+='__close'
            ;;
        d)
            cmd+='__delete'
            ;;
        s)
            cmd+='__sync'
            ;;
        q)
            cmd+='__quick'
            ;;
        h)
            cmd+='__help'
            ;;
        *)
            ;;
        esac
    done

    # Global options present in all commands
    opts='--header --color --csv --debug --namespace --indent \
    --project-namespace --help -h --version -v '

    case "${cmd}" in
    todoist)
        opts+='list l show completed-list c-l cl add a modify m close c \
        delete d labels projects karma sync s quick q help h'
        ;;

    todoist__add|todoist__modify)
        opts+='--priority -p --label-ids -L --project-id -P --project-name -N \
        --date -d --reminder -r'
        [ "$cmd" == 'todoist__modify' ] && opts+=' --content -c'

        case "${prev}" in
        --priority|-p)
            opts="1 2 3 4"
            ;;
        --label-ids|-L)
            COMPREPLY=( $(todoist labels | fzf --multi --select-1 --exit-0 \
            ${fzfquery} | cut -f 1 -d ' ' | paste -d, -s -) )
            return 0
            ;;
        --project-id|-P)
            COMPREPLY=( $(todoist projects | fzf --select-1 --exit-0 \
            ${fzfquery} | cut -f 1 -d ' ') )
            return 0
            ;;
        --project-name|-N)
            COMPREPLY=( "'$(todoist projects | fzf --select-1 --exit-0 \
            ${fzfquery} | cut -f 2- -d ' ' | cut -b 2- )'" )
            return 0
            ;;
        esac
        ;;

    todoist__list|todoist__completed-list)
        opts+='--filter -f'
        ;;

    todoist__show)
        opts+='--browse -o'
        ;;
    esac

    [ -n "$opts" ] && COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
}

complete -F _todoist todoist
