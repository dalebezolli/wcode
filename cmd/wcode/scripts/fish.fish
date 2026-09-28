function wcode
    if set -q argv[1]
        if not string match -q -- '-*' "$argv[1]"
            command wcode $argv
            return $status
        end
    end

    set -l _wcode_should_tmux true
    for _wcode_arg in $argv
        switch $_wcode_arg
            case -n --navigate-only
                set _wcode_should_tmux false
                break
        end
    end

    command wcode $argv
    set -l _wcode_status $status
    if test $_wcode_status -eq 3
        echo 'No project selected'
    end
    if test $_wcode_status -ne 0
        return $_wcode_status
    end

    set -l _wcode_dir (command wcode selection)
    set -l _wcode_status $status
    if test $_wcode_status -ne 0
        return $_wcode_status
    end
    set -l _wcode_name (string replace -r '^.*/' '' -- "$_wcode_dir")
    cd -- "$_wcode_dir"; or return $status

    if test "$_wcode_should_tmux" = false; or test -n "$TMUX"; or not type -q tmux
        return 0
    end

    if tmux has-session -t "$_wcode_name" 2>/dev/null
        tmux attach-session -t "$_wcode_name"
    else if test -f .tmux.conf
        tmux new-session -s "$_wcode_name" ';' source-file .tmux.conf
    else
        tmux new-session -s "$_wcode_name"
    end
end
