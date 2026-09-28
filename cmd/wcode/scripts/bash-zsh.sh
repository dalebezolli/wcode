wcode() {
  if [ "$#" -gt 0 ]; then
    case $1 in
      -*) ;;
      *) command wcode "$@"; return $? ;;
    esac
  fi

  local _wcode_arg _wcode_status _wcode_dir _wcode_name _wcode_should_tmux

  _wcode_should_tmux=true
  for _wcode_arg in "$@"; do
    case $_wcode_arg in
      -n|--navigate-only) _wcode_should_tmux=false; break ;;
    esac
  done

  command wcode "$@"
  _wcode_status=$?
  if [ "$_wcode_status" -eq 3 ]; then
    printf '%s\n' 'No project selected'
  fi
  if [ "$_wcode_status" -ne 0 ]; then
    return "$_wcode_status"
  fi

  _wcode_dir=$(command wcode selection) || return $?
  _wcode_name=${_wcode_dir##*/}
  cd -- "$_wcode_dir" || return

  if [ "$_wcode_should_tmux" = false ] || [ -n "${TMUX-}" ] || ! command -v tmux >/dev/null 2>&1; then
    return 0
  fi

  if tmux has-session -t "$_wcode_name" 2>/dev/null; then
    tmux attach-session -t "$_wcode_name"
  elif [ -f .tmux.conf ]; then
    tmux new-session -s "$_wcode_name" ';' source-file .tmux.conf
  else
    tmux new-session -s "$_wcode_name"
  fi
}
