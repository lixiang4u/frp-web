#!/bin/bash

VERSION=v0.1

PROJECT_PATH=$PWD

echo '开始打包web'
cd $PROJECT_PATH/frp-web-h5
npm i
npm run build

echo '开始打包可执行程序'
cd $PROJECT_PATH
go get
GOOS=linux   GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o frp-web-linux-amd64-$VERSION
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o frp-web-windows-amd64-$VERSION.exe
GOOS=darwin  GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o frp-web-darwin-amd64-$VERSION


cd $PROJECT_PATH
ls -lh  frp-web*

