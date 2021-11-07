#!/bin/bash
if [ "$1" = "" ]
then
    version=dev
else
    version=$1
fi
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=$version'" -o "discordcompressor-win64.exe" -v
GOOS=windows GOARCH=386 go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=$version'" -o "discordcompressor-win32.exe" -v
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=$version'" -o "discordcompressor-linux64" -v
GOOS=linux GOARCH=386 go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=$version'" -o "discordcompressor-linux32" -v