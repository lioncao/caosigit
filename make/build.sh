#!/bin/bash

build_dir(){
	cur_dir=$1
	cur_build_name=$2


	# echo "10 $0  11 $1  12 $2"

	cd $cur_dir

	if [ "x$cur_build_name" != "x" ]; then
		go_file_list=$(ls | grep *.go)
		if [ "$go_file_list"x != ""x ]; then			
			echo "go install $cur_build_name"
			go install $cur_build_name

		fi	
	fi


	dirList=$(ls -d)
	echo "$dirList"
	for d in $dirList; do
		if [ "x$d" != "x." ]; then			
			echo $d
			if [ "x$cur_build_name" != "x" ]; then
				next_build_name="$cur_build_name/$d"
			else
				next_build_name="$d"
			fi
			next_dir="$cur_dir/$d"
			build_dir $next_dir $next_build_name
		fi
		break
	done	
	return
}

cd ../src
# 源代码根目录
BUILD_NAME=""
SRC_DIR_FULL=$(pwd)
echo "00 $BUILD_NAME 01 $SRC_DIR_FULL"
build_dir $SRC_DIR_FULL $BUILD_NAME
