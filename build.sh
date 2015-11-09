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

#手动设置的默认的confd的配置文件路径
configDirFromBuild="$2"

if [ -n "$OS" ];then
   echo "use defined GOOS: "${OS}
else
   echo "use default GOOS: linux"
   OS=linux
fi

if [ -n "$configDirFromBuild" ];then
   echo "use defined configDirFromBuild: "${configDirFromBuild}
else
   echo "use default configDirFromBuild: linux"
   configDirFromBuild=""
fi

echo "start building with GOOS: "${OS}

export GOOS=${OS}
export GOARCH=amd64

if [ "$configDirFromBuild" != "" ];then
    flags="-X main.configDirFromBuild $configDirFromBuild"
    echo ${flags}
    go build -ldflags "$flags" -x -o confd confd.go
    go build -ldflags "$flags" -x -o confd-cli confd_cli.go
else
    go build -x -o confd confd.go
    go build -x -o confd-cli confd_cli.go
fi

cp confd docs/default_install_config/
cp confd-cli docs/default_install_config/

echo "finish building with GOOS: "${OS}" "${configDirFromBuild}