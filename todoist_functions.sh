select_items_command="todoist list | peco | cut -d ' ' -f 1 | tr '\n' ' '"

# todoist find item
function peco-todoist-item () {
    local SELECTED_ITEMS="$(eval ${select_items_command})"
    if [ -n "$SELECTED_ITEMS" ]; then
        if [ -n "$LBUFFER" ]; then
            local new_left="${LBUFFER%\ } $SELECTED_ITEM"
        else
            local new_left="$SELECTED_ITEMS"
        fi
        BUFFER=${new_left}${RBUFFER}
        CURSOR=${#new_left}
    fi
}
zle -N peco-todoist-item
bindkey "^xtt" peco-todoist-item

# todoist find project
function peco-todoist-project () {
    local SELECTED_PROJECT="$(todoist projects | peco | head -n1 | cut -d ' ' -f 1)"
    if [ -n "$SELECTED_PROJECT" ]; then
        if [ -n "$LBUFFER" ]; then
            local new_left="${LBUFFER%\ } -P $SELECTED_PROJECT"
        else
            local new_left="-P $SELECTED_PROJECT"
        fi
        BUFFER=${new_left}${RBUFFER}
        CURSOR=${#new_left}
    fi
}
zle -N peco-todoist-project
bindkey "^xtp" peco-todoist-project

# todoist find labels
function peco-todoist-labels () {
    local SELECTED_LABELS="$(todoist labels | peco | cut -d ' ' -f 1 | tr '\n' ',' | sed -e 's/,$//')"
    if [ -n "$SELECTED_LABELS" ]; then
        if [ -n "$LBUFFER" ]; then
            local new_left="${LBUFFER%\ } -L $SELECTED_LABELS"
        else
            local new_left="-L $SELECTED_LABELS"
        fi
        BUFFER=${new_left}${RBUFFER}
        CURSOR=${#new_left}
    fi
}
zle -N peco-todoist-labels
bindkey "^xtl" peco-todoist-labels

# todoist close
function peco-todoist-close() {
    local SELECTED_ITEMS="$(eval ${select_items_command})"
    if [ -n "$SELECTED_ITEMS" ]; then
        BUFFER="todoist close $(echo "$SELECTED_ITEMS" | tr '\n' ' ')"
        CURSOR=$#BUFFER
    fi
    zle accept-line
}
zle -N peco-todoist-close
bindkey "^xtc" peco-todoist-close

# todoist delete
function peco-todoist-delete() {
    local SELECTED_ITEMS="$(eval ${select_items_command})"
    if [ -n "$SELECTED_ITEMS" ]; then
        BUFFER="todoist delete $(echo "$SELECTED_ITEMS" | tr '\n' ' ')"
        CURSOR=$#BUFFER
    fi
    zle accept-line
}
zle -N peco-todoist-delete
bindkey "^xtk" peco-todoist-delete
