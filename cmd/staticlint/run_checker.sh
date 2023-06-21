#!/bin/sh

go build ./cmd/staticlint/

./staticlint -test=false ./...