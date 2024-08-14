#!/bin/bash

tmux-workspace "NWGo-CLI" "editor" -c "nvim && zsh" \
    -w "nwgo-demo" -c "cd ../nwgo-demo && zsh"
