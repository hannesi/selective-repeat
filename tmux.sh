#!/bin/sh

# uses 1 based indexing

session_name="go-back-n"

tmux has-session -t $session_name 2>/dev/null

if [ $? != 0 ]; then
  tmux new-session -d -s $session_name
  tmux new-window -t $session_name -d
  tmux split-window -h -t $session_name:2.1
fi

tmux attach -t $session_name

