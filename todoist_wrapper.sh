## TODOIST
todoist sync
##BASIC FUNCS##
function filter () {##prints all things according to a given filter
	todoist list --filter $1 
}
function t-get-tasks() {##opens peco on global task list and returns task ids that are selected
  TASK_LIST=$(eval todoist --namespace --project-namespace list | peco | cut -d ' ' -f 1 | tr '\n' ' ')
  echo "$TASK_LIST"
}
function t-get-proj() {##gets the project id or name from the project-namespace
	PROJ="$(todoist --project-namespace projects | peco --prompt $1 | head -n1 | cut -d ' ' -f $2)"
  echo "$PROJ"
}
function t-get-tasks-from-proj() {##prompts to select a project first, then prompts to select tasks from the project
	TASK_LIST="$(eval t-prt-proj | peco --prompt "tasks>"| cut -d ' ' -f 1 | tr '\n' ' ')"
	echo "$TASK_LIST"
}
function t-get-tasks-from-filter() {##prompts to select a project first, then prompts to select tasks from the project
	TASK_LIST="$(eval filter "$1" | peco --prompt "$1>"| cut -d ' ' -f 1 | tr '\n' ' ')"
	echo "$TASK_LIST"
}
function t-prt-tasks() {##prints generic list of task ids
  for TASK in $*
  do
    todoist show --browse "$TASK"
  done
}
function t-prt-proj() {##prints a project selected with peco
	PROJ="$(eval "t-get-proj 'projects>' 2")"
	filter "$PROJ"
}
function t-add-task-w-proj() {##prompts for desired project and adds the passed in string as a task in it
##come_back
	PROJ="$(eval "t-get-proj 'projects>' 1")"
	todoist add $1 -P "$PROJ"
}
function t-mv-tasks() {##recieves list of tasks, prompts for a destination project and moves them there
  TARGET_PROJECT=$(eval "t-get-proj 'Pick ONE destination project>' 1")
  if [ -n "$TARGET_PROJECT" ]; then
    for TASK in $*
    do
      todoist m "$TASK" -P "$TARGET_PROJECT"
      t-prt-tasks "$TASK"
    done
  fi
}
function y-n-valid() {
  if [ "$1" = "y" -o "$1" = "Y" ]; then
    echo "1"
  elif [ "$1" = "n" -o "$1" = "N" ]; then
    echo "0"
  else
    echo "2"
  fi
}
function t-del-tasks () {
	if [ -n "$1" ]; then
    choice="-1"
    while ! [ "$choice" = "1" -o "$choice" = "0" ]
    do
      read bulk\?"Bulk delete? [y/n]"
      choice=$(y-n-valid "$bulk")
    done

    if [[ "$choice" = "1" ]]; then
      t-prt-tasks $*
      choice="-1"
      while ! [ "$choice" = "1" -o "$choice" = "0" ]
      do
        read delete\?"Are you sure you want to delete these tasks?"
        choice=$(y-n-valid "$delete")
      done
      if [[ "$choice" = "1" ]]; then
        for TASK in $*
        do
          todoist d "$TASK"
        done
      else
        echo "Exiting tdel."
        return 0
      fi
    else
	    for TASK in $*
      do
        choice="-1"
        t-prt-tasks "$TASK"
        while ! [ "$choice" = "1" -o "$choice" = "0" ]
        do
          read delete\?"Delete this task?"
          choice=$(y-n-valid "$delete")
        done

        if [[ "$choice" = "1" ]]; then
          todoist d $TASK
        else
          echo "Ok. Skipping to next task."
        fi
      done
    fi
  else
    echo "Recieved empty task list in t-del-tasks"
	fi
}
function t-mv-today() {
  for TASK in $*
  do
    t-prt-tasks "$TASK"
    todoist m "$TASK" -d "today"
  done
}
function t-rm-today() {
  for TASK in $*
  do
    todoist m "$TASK" -d 'no date'
  done
}
##todoist_aliases
##gets task list and prints each task selected
alias tsk='t-prt-tasks $(eval t-get-tasks)'
#gets task list from global task list and prompts to delete them
alias tdelg='t-del-tasks $(eval t-get-tasks)'
##gets task list from global task list
alias tmvg='t-mv-tasks $(eval t-get-tasks)'
##gets task list from individual project list 
alias tmv='t-mv-tasks $(eval t-get-tasks-from-proj)'
#gets task list from individual project list and prompts to delete them
alias tdel='t-del-tasks $(eval t-get-tasks-from-proj)'
#gets task list from individual project list and moves them to today's todo list
alias mvtdy='t-mv-today $(eval t-get-tasks-from-proj)'
#gets task list from individual project list and removes them to today's todo list
alias rmtdy='t-rm-today $(eval "t-get-tasks-from-filter 'today'")'

#shortened command call
alias td='todoist'
#syncs todoist
alias tsync='todoist sync'
#prints the Inbox
alias tin='filter "#Inbox"'
#prints today's todo list
alias tdy='filter "today"'
#prints this week's tasks

#add task to a project chosen through peco
alias ta='t-add-task-w-proj'
#print a project chosen through peco
alias tp='t-prt-proj'
#alias for todoist_functions.sh widget which completes a task
#temporary until I can get selection working through project and today selection
alias tfin='xdotool key "control+x" "t" "c"'