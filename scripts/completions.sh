#!/bin/sh
set -e
rm -rf completions
mkdir completions
go build -o tt .
for sh in bash zsh fish; do
	./tt completion "$sh" >"completions/tt.$sh"
done
