#!/usr/bin/env bash

covermode=${COVERMODE:-atomic}
coverdir=$(mktemp -d /tmp/coverage.XXXXXXXXXX)
profile="${coverdir}/cover.out"

echo "profile: ${profile}"

go test -race -coverprofile=${profile} -covermode="${covermode}" ./...
go tool cover -func ${profile}
