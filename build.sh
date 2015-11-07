#!/bin/sh

#####################################################################
# usage:
# sh build.sh 构建默认的linux64位程序
# sh build.sh darwin(或linux), 构建指定平台的64为程序

# examples:
# sh build.sh darwin 构建MacOS版本的程序
# sh build.sh linux 构建linux版本的程序
#####################################################################

source /etc/profile

OS="$1"

if [ -n "$OS" ];then
   echo "use defined GOOS: "${OS}
else
   echo "use default GOOS: linux"
   OS=linux
fi

echo "start building with GOOS: "${OS}

export GOOS=${OS}
export GOARCH=amd64


flags="-X main.buildstamp `date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash `git rev-parse HEAD`"
echo ${flags}
go build -ldflags "$flags" -x -o docs/default_install_config/confd confd.go


echo "finish building with GOOS: "${OS}