#!/usr/bin/env bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ ! -e "$CURRENT_DIR/bin/hazcld" ]]; then
    pushd "$CURRENT_DIR" >/dev/null 2>&1 || exit $?
    command go build -o "$CURRENT_DIR"/bin/hazcld
    popd >/dev/null 2>&1 || exit $?
fi

# Below code is copied and adapted from vim-tmux-navigator
# See: https://github.com/christoomey/vim-tmux-navigator/blob/6a1e58c3ca3bc7acca36c90521b3dfae83b2a602/vim-tmux-navigator.tmux#L1

version_pat='s/^tmux[^0-9]*([.0-9]+).*/\1/p'

is_vim="$CURRENT_DIR/bin/hazcld 'g?(view|n?vim?x?)(diff)?' '#{pane_pid}'"
echo "$is_vim"

tmux bind-key -n C-h if-shell "$is_vim" "send-keys C-h" "select-pane -L"
tmux bind-key -n C-j if-shell "$is_vim" "send-keys C-j" "select-pane -D"
tmux bind-key -n C-k if-shell "$is_vim" "send-keys C-k" "select-pane -U"
tmux bind-key -n C-l if-shell "$is_vim" "send-keys C-l" "select-pane -R"
tmux_version="$(tmux -V | sed -En "$version_pat")"
tmux setenv -g tmux_version "$tmux_version"

tmux if-shell -b '[ "$(echo "$tmux_version < 3.0" | bc)" = 1 ]' \
    "bind-key -n 'C-\\' if-shell \"$is_vim\" 'send-keys C-\\'  'select-pane -l'"
tmux if-shell -b '[ "$(echo "$tmux_version >= 3.0" | bc)" = 1 ]' \
    "bind-key -n 'C-\\' if-shell \"$is_vim\" 'send-keys C-\\\\'  'select-pane -l'"

tmux bind-key -T copy-mode-vi C-h select-pane -L
tmux bind-key -T copy-mode-vi C-j select-pane -D
tmux bind-key -T copy-mode-vi C-k select-pane -U
tmux bind-key -T copy-mode-vi C-l select-pane -R
tmux bind-key -T copy-mode-vi C-\\ select-pane -l
