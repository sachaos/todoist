bind '"\C-xtc": "todoist list | fzf | cut -d \" \" -f 1 | xargs -I {} todoist close "{}"\n" '
bind '"\C-xtk": "todoist list | fzf | cut -d \" \" -f 1 | xargs -I {} todoist delete "{}"\n" '
bind '"\C-xtm": "todoist list | fzf | cut -d \" \" -f 1 | xargs -I {} todoist modify \"{}\" -c " '
bind '"\C-xtl": "todoist list | fzf | cut -d \" \" -f 1 | xargs -I {} todoist show "{}"\n"'
