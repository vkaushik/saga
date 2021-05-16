#!/bin/bash

if [ -n "$(gofmt -l . | grep -v vendor | grep -v .gopath)" ]; then
  echo "Go code is not formatted:"
  for i in $(gofmt -l . | grep -v vendor | grep -v .gopath); do
    echo "Fix this file: $i"
    gofmt -d "$i"
  done
  exit 1
fi