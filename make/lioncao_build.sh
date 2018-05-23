#!/bin/bash
list=""
# list="$list lioncao"
# list="$list lioncao/net"
list="$list lioncao/net/http"
list="$list lioncao/net/socket"
list="$list lioncao/net/websocket"
# list="$list lioncao/util"
list="$list lioncao/util/cmd"
# list="$list lioncao/util/db"
list="$list lioncao/util/db/mongodb"
list="$list lioncao/util/db/redis"
list="$list lioncao/util/network"
list="$list lioncao/util/service"
list="$list lioncao/util/tools"
list="$list lioncao/util/msgcode"


echo "=========================================="
echo "=========================================="
echo "GOROOT=$GOROOT"
echo "GOPATH=$GOPATH"
echo "=========================================="
echo "=========================================="




for name in $list; do
	go install $name
	if [ $? -ne 0 ]; then
		echo "=========================================="
		echo "install failed $name"
		echo "=========================================="
		# exit 0
	else
		echo "install ok $name"
	fi
done

