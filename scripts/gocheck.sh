#!/usr/bin/env bash

GO_FILES=$(find . -iname '*.go' -type f | grep -v /vendor/) # All the .go files, excluding vendor/

# Check golint
if [[ $(golint ${GO_FILES}) ]]; then
    echo 'Some files need linting attention.'
    echo 'You can use the command: \`golint .\` to find out more info.'
    exit 1
fi

# Check gofmt
if [[ $(gofmt -s -l ${GO_FILES}) ]]; then
    echo 'gofmt needs running on some files.'
    echo 'You can use the command: \`gofmt -w .\` to reformat code.'
    exit 1
fi

# Check go vet
if [[ $(go vet ${GO_FILES}) ]]; then
    echo 'Some files need attention.'
    echo 'You can use the command: \`go vet .\` to find out more info.'
    exit 1
fi

exit 0
