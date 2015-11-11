#!/usr/bin/env bash

#confd config dir
DEFAULT_PATH="$1"

if [ -n "$DEFAULT_PATH" ];then
   echo "use defined config dir: "${DEFAULT_PATH}
else
   echo "use default config dir: /data/confd"
   DEFAULT_PATH=/data/confd
fi

mkdir -p $DEFAULT_PATH

#make a backup for the last configurations
CURRENT_TIME=$(date +%Y%m%d-%H%M%S)
tar zcvf $DEFAULT_PATH".${CURRENT_TIME}.tar.gz" $DEFAULT_PATH
rm -rf $DEFAULT_PATH/*


echo "Install confd in default path: $DEFAULT_PATH"

mkdir -p $DEFAULT_PATH/data
mkdir -p $DEFAULT_PATH/meta
mkdir -p $DEFAULT_PATH/templates


cp ./files/config.toml $DEFAULT_PATH/data/
cp ./files/filestore.toml $DEFAULT_PATH/data/
cp ./files/example_1.toml $DEFAULT_PATH/meta/
cp ./files/example_2.toml $DEFAULT_PATH/meta/
cp ./files/example.tmpl $DEFAULT_PATH/templates/

cp ./files/confd $DEFAULT_PATH/
cp ./files/confd-cli $DEFAULT_PATH/

echo "Installed confd"
echo "By default, confd uses a file store, you also can choose redis or zookeeper"
echo 'Now you can:'
echo "  1. cd $DEFAULT_PATH"
echo '  2. use "./confd --confdir='$DEFAULT_PATH'" to see how confd works.'
echo "  3. check folders[data/meta/templates] under $DEFAULT_PATH to learn how to use confd"