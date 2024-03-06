#!/bin/bash

PROJECT_PATH=$PWD

echo '开始打包web'
cd $PROJECT_PATH/frp-web-h5
npm i
npm run build

echo '开始打包可执行程序'
cd $PROJECT_PATH
go get
GOOS=linux GOARCH=amd64 go build -o frp-web-linux-amd64-v0.1
GOOS=windows GOARCH=amd64 go build -o frp-web-windows-amd64-v0.1.exe
GOOS=darwin GOARCH=amd64 go build -o frp-web-darwin-amd64-v0.1.app

ls -lh  frp-web*

