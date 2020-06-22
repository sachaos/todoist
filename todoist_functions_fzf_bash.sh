#!/usr/bin/env bash
# Bash completion for https://github.com/sachaos/todoist
# Thanks to other completion scripts (like rustup and gh) for inspiration.

__todoist_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

_todoist() {
    local arg cur prev opts cmd fzfquery fzfcmd fzftasks
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts=''
    cmd=''
    fzfquery=''

    [ -n "$cur" ] && fzfquery="-q $(echo $cur | tr -d " " | tr -d "'")"
    fzfcmd="fzf --select-1 --exit-0 $fzfquery"
    fzftasks=0

    __todoist_debug "${FUNCNAME[0]}(): cur=$cur / prev=$prev"\
" / fzfquery=$fzfquery / COMP_WORDS[@]=${COMP_WORDS[@]}"

    for arg in "${COMP_WORDS[@]}"; do
        case "$arg" in
        todoist)
            cmd='todoist'
            ;;
        # These are the current commands; not all have completion options,
        # but they're listed here anyway, for the future
        list|show|completed-list|add|modify|close|delete|labels|projects|\
        karma|sync|quick|help)
            [ "$cmd" == 'todoist' ] && cmd+="__$arg"
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
    __todoist_debug "${FUNCNAME[0]}(): cmd=$cmd"

    # Global options present in all commands
    opts='--header --color --csv --debug --namespace --indent'\
' --project-namespace --help -h --version -v '

    case "$cmd" in
    todoist)
        opts+='list l show completed-list c-l cl add a modify m close c'\
' delete d labels projects karma sync s quick q help h'
        ;;

    todoist__add|todoist__modify)
        opts+='--priority -p --label-ids -L --project-id -P --project-name -N'\
' --date -d'

        case "$cmd" in
        todoist__add)
            opts+=' --reminder -r'
            ;;
        todoist__modify)
            opts+=' --content -c'
            ;;
        esac

        case "$prev" in
        --priority|-p)
            opts="1 2 3 4"
            ;;
        --label-ids|-L)
            # shellcheck disable=SC2207
            COMPREPLY=( $(todoist labels | $fzfcmd --multi | cut -d ' ' -f 1 \
            | paste -d, -s -) )
            return 0
            ;;
        --project-id|-P)
            # shellcheck disable=SC2207
            COMPREPLY=( $(todoist projects | $fzfcmd | cut -d ' ' -f 1) )
            return 0
            ;;
        --project-name|-N)
            COMPREPLY=( "'$(todoist projects | $fzfcmd | cut -d ' ' -f 2- \
            | cut -b 2- )'" )
            return 0
            ;;
        *)
            __todoist_debug "cmd=$cmd / cur=$cur / cur:0:1=${cur:0:1}"
            if [ "$cmd" == 'todoist__modify' ] && [ "${cur:0:1}" != '-' ]; then
                # If it's not an option, list tasks
                fzftasks=1
            fi
            ;;
        esac
        ;;

    todoist__list|todoist__completed-list)
        opts+='--filter -f'
        ;;

    todoist__show)
        opts+='--browse -o'
        # If it's not an option, list tasks
        [ "${cur:0:1}" != '-' ] && fzftasks=1
        ;;

    todoist__close|todoist__delete)
        # If it's not an option, list tasks
        __todoist_debug "cmd=$cmd / cur=$cur / cur:0:1=${cur:0:1}"
        [ "${cur:0:1}" != '-' ] && fzftasks=1
        ;;
    esac


    if [ $fzftasks -eq 1 ]; then
        __todoist_debug "fzfcmd=$fzfcmd"

        # shellcheck disable=SC2207
        COMPREPLY=( $(todoist --namespace --project-namespace list \
        | $fzfcmd | cut -d ' ' -f 1 | tr -d "'") )
        return 0
    fi

    __todoist_debug "${FUNCNAME[0]}(): opts=$opts"

    # shellcheck disable=SC2207
    [ -n "$opts" ] && COMPREPLY=( $(compgen -W "$opts" -- "$cur") )
}

complete -F _todoist todoist
