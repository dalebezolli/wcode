wcode() {
  if [[ $# -gt 0 && $1 != -* ]]; then
    command wcode "$@"
    return $?
  fi

  local should_run_tmux=true
  local arg
  for arg in "$@"; do
    if [[ $arg == -n || $arg == --navigate-only ]]; then
      should_run_tmux=false
      break
    fi
  done

  command wcode "$@"
  local status=$?
  if (( status == 3 )); then
    echo "No project selected"
  fi
  if (( status != 0 )); then
    return "$status"
  fi

  local working_dir name
  working_dir=$(<"$HOME/.config/wcode/selection") || return 2
  name=${working_dir##*/}
  cd -- "$working_dir" || return

  if ! command -v tmux >/dev/null 2>&1 || [[ -n ${TMUX-} || $should_run_tmux != true ]]; then
    return 0
  fi

  local -a args
  if tmux has-session -t "$name" 2>/dev/null; then
    args=(attach-session -t "$name")
  else
    args=(new-session -s "$name")
    if [[ -f .tmux.conf ]]; then
      args+=(';' source-file .tmux.conf)
    fi
  fi

  tmux "${args[@]}"
}
