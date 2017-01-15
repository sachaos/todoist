select_items_command="todoist list | peco | cut -d ' ' -f 1 | tr '\n' ' '"

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
function peco-todoist-item () {
    local SELECTED_ITEMS="$(eval ${select_items_command})"
    insert-in-buffer $SELECTED_ITEMS
}
zle -N peco-todoist-item
bindkey "^xtt" peco-todoist-item

# todoist find project
function peco-todoist-project () {
    local SELECTED_PROJECT="$(todoist projects | peco | head -n1 | cut -d ' ' -f 1)"
    insert-in-buffer "${SELECTED_PROJECT}" "-P"
}
zle -N peco-todoist-project
bindkey "^xtp" peco-todoist-project

# todoist find labels
function peco-todoist-labels () {
    local SELECTED_LABELS="$(todoist labels | peco | cut -d ' ' -f 1 | tr '\n' ',' | sed -e 's/,$//')"
    insert-in-buffer "${SELECTED_LABELS}" "-L"
}
zle -N peco-todoist-labels
bindkey "^xtl" peco-todoist-labels

function todoist-exec-with-select-task () {
    if [ -n "$2" ]; then
        BUFFER="todoist $1 $(echo "$2" | tr '\n' ' ')"
        CURSOR=$#BUFFER
        zle accept-line
    fi
}

# todoist close
function peco-todoist-close() {
    local SELECTED_ITEMS="$(eval ${select_items_command})"
    todoist-exec-with-select-task close $SELECTED_ITEMS
}
zle -N peco-todoist-close
bindkey "^xtc" peco-todoist-close

# todoist delete
function peco-todoist-delete() {
    local SELECTED_ITEMS="$(eval ${select_items_command})"
    todoist-exec-with-select-task delete $SELECTED_ITEMS
}
zle -N peco-todoist-delete
bindkey "^xtk" peco-todoist-delete
