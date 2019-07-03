select_items_command="todoist --namespace --project-namespace list | fzf | cut -d ' ' -f 1 | tr '\n' ' '"

function insert-in-buffer () {
    if [ -n "$1" ]; then
        local new_left=""
        if [ -n "$LBUFFER" ]; then
            new_left="${new_left}${LBUFFER} "
        fi
        if [ -n "$2" ]; then
            new_left="${new_left}${2} "
        fi
        new_left="${new_left}$1"
        BUFFER=${new_left}${RBUFFER}
        CURSOR=${#new_left}
    fi
}

# todoist find item
function fzf-todoist-item () {
    local SELECTED_ITEMS="$(eval ${select_items_command})"
    insert-in-buffer $SELECTED_ITEMS
}
zle -N fzf-todoist-item
bindkey "^xtt" fzf-todoist-item

# todoist find project
function fzf-todoist-project () {
    local SELECTED_PROJECT="$(todoist --project-namespace projects | fzf | head -n1 | cut -d ' ' -f 1)"
    insert-in-buffer "${SELECTED_PROJECT}" "-P"
}
zle -N fzf-todoist-project
bindkey "^xtp" fzf-todoist-project

# todoist find labels
function fzf-todoist-labels () {
    local SELECTED_LABELS="$(todoist labels | fzf | cut -d ' ' -f 1 | tr '\n' ',' | sed -e 's/,$//')"
    insert-in-buffer "${SELECTED_LABELS}" "-L"
}
zle -N fzf-todoist-labels
bindkey "^xtl" fzf-todoist-labels

# todoist select date
function fzf-todoist-date () {
    date -v 1d &>/dev/null
    if [ $? -eq 0 ]; then
        # BSD date option
        OPTION="-v+#d"
    else
        # GNU date option
        OPTION="-d # day"
    fi

    local SELECTED_DATE="$(seq 0 30 | xargs -I# date $OPTION '+%d/%m/%Y %a' | fzf | cut -d ' ' -f 1)"
    insert-in-buffer "'${SELECTED_DATE}'" "-d"
}
zle -N fzf-todoist-date
bindkey "^xtd" fzf-todoist-date

function todoist-exec-with-select-task () {
    if [ -n "$2" ]; then
        BUFFER="todoist $1 $(echo "$2" | tr '\n' ' ')"
        CURSOR=$#BUFFER
        zle accept-line
    fi
}

# todoist close
function fzf-todoist-close() {
    local SELECTED_ITEMS="$(eval ${select_items_command})"
    todoist-exec-with-select-task close $SELECTED_ITEMS
}
zle -N fzf-todoist-close
bindkey "^xtc" fzf-todoist-close

# todoist delete
function fzf-todoist-delete() {
    local SELECTED_ITEMS="$(eval ${select_items_command})"
    todoist-exec-with-select-task delete $SELECTED_ITEMS
}
zle -N fzf-todoist-delete
bindkey "^xtk" fzf-todoist-delete

# todoist open
function fzf-todoist-open() {
    local SELECTED_ITEMS="$(eval ${select_items_command})"
    todoist-exec-with-select-task "show --browse" $SELECTED_ITEMS
}
zle -N fzf-todoist-open
bindkey "^xto" fzf-todoist-open
