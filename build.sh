#!/bin/sh

mod='github.com/SongCastle/ggnb'
app="/go/src/${mod}"
main="${app}/main.go"

cd $app

if [ ! -e "${app}/go.mod" ]; then
  go mod init $mod
fi

if [ ! -e $main ]; then
  echo 'main does not exist'
  exit 1
fi

bin='/bin/ggnb'

# 参考
# https://zenn.dev/spiegel/articles/20210223-go-module-aware-mode#go-mod-tidy-%E3%81%AB%E3%82%88%E3%82%8B%E3%83%A2%E3%82%B8%E3%83%A5%E3%83%BC%E3%83%AB%E6%83%85%E5%A0%B1%E3%81%AE%E6%9B%B4%E6%96%B0
go mod tidy
GOOS=linux CGO_ENABLED=0 go build -o $bin

if [ ! $? = 0 ]; then
  echo 'build filed'
  exit 1
fi

$bin
