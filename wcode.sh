#!/bin/bash

SCRIPT_DIR=$(cd $(dirname $BASH_SOURCE) && pwd)

wcode() {
  $SCRIPT_DIR/bin/wcode;
  if [ $(echo $?) -ne 0 ]; then

    echo "No project selected";
  fi

  WORKING_DIR=$(cat ~/.config/wcode/selection)
  NAME=$(echo $WORKING_DIR | grep -o "[^/]*$")

  cd $WORKING_DIR;

  which tmux 2>/dev/null 1>/dev/null
  HAS_TMUX=$?

  # If tmux doesn't exist or is inside an existing tmux session do an early return
  if [ $HAS_TMUX -eq 1 -o -n "$TMUX" ]; then
    return
  fi

  tmux has-session -t $NAME 2> /dev/null
  HAS_SESSION=$?

  ARGS=()
  if [ $HAS_SESSION -eq 1 ]; then
    ARGS+=(new-session -s $NAME)
  else
    ARGS+=(attach-session -t $NAME)

  fi

  if [ $HAS_SESSION -eq 1 -a -e ".tmux.conf" ]; then
    ARGS+=(';' source-file .tmux.conf)
  fi

  tmux "${ARGS[@]}"
}
