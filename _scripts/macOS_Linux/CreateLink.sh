#!/bin/bash
file_path=$(cd "$(dirname "$0")"; pwd)
project_path=$file_path"/../../../EasyGA"

cd ~

read -r -p "Force mode? [Y/n] " input
case $input in
    [yY][eE][sS]|[yY])
		echo "Force Mode"
		if [ ! -d "go/src/github.com/fang2hou" ]; then
			mkdir ~/go/src/github.com/fang2hou
		fi
		ln -snf $project_path ~/go/src/github.com/fang2hou/easyga
		;;

    [nN][oO]|[nN])
		echo "Normal Mode"
		if [ ! -d "go/src/github.com/fang2hou" ]; then
			mkdir ~/go/src/github.com/fang2hou
		fi
		ln -sn $project_path ~/go/src/github.com/fang2hou/easyga
       	;;

    *)
	echo "Invalid input..."
	exit 1
	;;
esac
