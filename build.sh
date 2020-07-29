#!/bin/sh

echo "Cleaning dist directory..."
rm -rf dist/*
echo "Making dist directory..."
mkdir dist
echo "Building darwin amd64..."
GOOS=darwin GOARCH=amd64 go build -o dist/filesSearch_darwin-amd64.out -ldflags "-s" filesSearch.go
echo "Building windows amd64..."
GOOS=windows GOARCH=amd64 go build -o dist/filesSearch_win-amd64.exe -ldflags "-s" filesSearch.go
echo "Building windows x86..."
GOOS=windows GOARCH=386 go build -o dist/filesSearch_win-x86.exe -ldflags "-s" filesSearch.go
echo "Building linux amd64..."
GOOS=linux GOARCH=amd64 go build -o dist/filesSearch_linux-amd64.out -ldflags "-s" filesSearch.go

echo "Copying search.properties into dist..."
cp search.properties dist/

echo "Completed"