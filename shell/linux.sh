#!/bin/sh

COMMIT_ID=$(git log | head -n 1 | awk '{print $2;}')
AUTHOR=$(git log | head -n 3 | grep Author | awk '{print $2;}')
BRANCH_NAME=$(git branch | awk '/\*/ { print $2; }')
GIT_INFO="$COMMIT_ID,$AUTHOR,$BRANCH_NAME"
echo "$GIT_INFO"

go mod tidy
go mod vendor
GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/cs-sea/rabbit/compile.GitInfo=$GIT_INFO -w -s" -v -o rabbit main/main.go
