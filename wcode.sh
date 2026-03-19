#!/bin/bash

SCRIPT_DIR=$(cd $(dirname $BASH_SOURCE) && pwd)

# Exit statuses
EXIT_OK=0
EXIT_NO_PROJECTS=1
EXIT_BAD_PATH=2
EXIT_NO_SELECTION=3
EXIT_TERMINATED=9

wcode() {
  SHOULD_RUN_TMUX="true"
  for arg in "$@"; do
    if [[ "$arg" == "-n" || "$arg" == "--navigate-only" ]]; then
      SHOULD_RUN_TMUX="false"
      break
    fi
  done

  $SCRIPT_DIR/bin/wcode $@;
  WCODE_STATUS=$?

  if [ $WCODE_STATUS -eq $EXIT_NO_SELECTION ]; then
    echo "No project selected";
  fi

  if [ $WCODE_STATUS -ne $EXIT_OK ]; then
    return;
  fi

  WORKING_DIR=$(cat ~/.config/wcode/selection)
  NAME=$(echo $WORKING_DIR | grep -o "[^/]*$")

  cd $WORKING_DIR;

  which tmux 2>/dev/null 1>/dev/null
  HAS_TMUX=$?

  # If tmux doesn't exist or is inside an existing tmux session do an early return
  if [[ $HAS_TMUX -eq 1 || -n "$TMUX" || "$SHOULD_RUN_TMUX" != "true" ]]; then
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
