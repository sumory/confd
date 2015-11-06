#!/usr/bin/env bash

DEFAULT_PATH=/data/confd
mkdir -p $DEFAULT_PATH
echo "Install confd in default path: $DEFAULT_PATH"

mkdir -p $DEFAULT_PATH/data
mkdir -p $DEFAULT_PATH/meta
mkdir -p $DEFAULT_PATH/templates


cp ./default_install_config/config.toml $DEFAULT_PATH/data/
cp ./default_install_config/filestore.toml $DEFAULT_PATH/data/
cp ./default_install_config/example_1.toml $DEFAULT_PATH/meta/
cp ./default_install_config/example_2.toml $DEFAULT_PATH/meta/
cp ./default_install_config/example.tmpl $DEFAULT_PATH/templates/

cp ./default_install_config/confd $DEFAULT_PATH/

echo "Installed confd"
echo "By default, confd uses a file store, you also can choose redis or zookeeper"
echo 'Now you can:'
echo "  1. cd $DEFAULT_PATH"
echo '  2. use "./confd --confdir='$DEFAULT_PATH'" to see how confd works.'
echo "  3. check folders[data/meta/templates] under $DEFAULT_PATH to learn how to use confd"