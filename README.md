# hazcld - Check if a process has a certain child

`hazcld` is a Linux utility that given a process ID checks if the
corresponding process has a child process with the name matching the
passed regex.

Some tools require this information to determine how to process user
interaction.
[`vim-tmux-navigator`](https://github.com/christoomey/vim-tmux-navigator)
uses the
[following](https://github.com/christoomey/vim-tmux-navigator/blob/6a1e58c3ca3bc7acca36c90521b3dfae83b2a602/vim-tmux-navigator.tmux#L5)
to determine if a tmux pane runs `vim`:

    is_vim="ps -o state= -o comm= -t '#{pane_tty}' \
        | grep -iqE '^[^TXZ ]+ +(\\S+\\/)?g?(view|n?vim?x?)(diff)?$'"

This however breaks if the pane executes `vim` inside a `poetry shell
session`.

The following shell script would to the trick, but is rather slow:

```bash
#!/usr/bin/env bash

PATTERN="[[:space:]]*g?(view|n?vim?x?)(diff)?"

function find_vim() {
    local pid

    pid="$1"

    while [[ -n "$pid" ]]; do
        read -r pid cmd <<<"$(command ps --ppid "$pid" -o pid= -o comm=)"
        if [[ "$cmd" =~ $PATTERN ]]; then
            return 0
        fi
    done

    return 1
}

find_vim "$1"
```

This is an attempt to generalize the above script and make it faster by
reducing the amount of calls to `ps`.

A short try with the initial draft version looks promising. Please note
that this was not a properly planned benchmark. It mainly served to play
around a little with
[`hyperfine`](https://github.com/sharkdp/hyperfine):

| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `./scripts/tmux-pane-runs-vim 2210` | 22.6 ± 1.2 | 21.1 | 26.9 | 5.18 ± 0.76 |
| `./bin/hazcld '\s*g?(view\|n?vim?x?)(diff)?' 2210` | 4.4 ± 0.6 | 3.6 | 6.6 | 1.00 |

Above table shows the output of

    hyperfine --export-markdown results.md \
              --warmup 10 \
              "./scripts/tmux-pane-runs-vim 2210" \
              "./bin/hazcld '\s*g?(view|n?vim?x?)(diff)?' 2210"

## Installation

The package can be installed using `go get`:

    go get github.com/fhofherr/hazcld

## Usage

Call `hazcld` using a [Go regular
expression](https://pkg.go.dev/regexp/syntax) matching the command the
child process you are interested in is executing and the process ID of a
parent process.

    hazcld '\s*g?(view|n?vim?x?)(diff)?' 2210

### Tmux

This repository contains a `hazcld_vim.tmux` script which can be used as
a [TPM](https://github.com/tmux-plugins/tpm) tmux plugin.

Add the following line to your `tmux.conf` and install the
`christoomey/vim-tmux-navigator` vim plugin.

    set -g @plugin 'fhofherr/hazcld'

## License

Copyright © 2021 Ferdinand Hofherr

Distributed under the MIT License.
