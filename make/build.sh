#!/bin/bash
list=""
# list="$list 3rdparty"
# list="$list 3rdparty/goArrayList"
list="$list 3rdparty/goArrayList/goArrayList"
list="$list 3rdparty/mahonia"
list="$list 3rdparty/ssdb"
# list="$list buddy"
# list="$list buddy/net"
list="$list buddy/net/http"
list="$list buddy/net/socket"
# list="$list buddy/util"
list="$list buddy/util/cmd"
# list="$list buddy/util/db"
list="$list buddy/util/db/mongodb"
list="$list buddy/util/db/redis"
list="$list buddy/util/network"
list="$list buddy/util/service"
list="$list buddy/util/tools"
# list="$list code.google.com"
# list="$list code.google.com/p"
# list="$list code.google.com/p/goprotobuf"
# list="$list code.google.com/p/goprotobuf/lib"
# list="$list code.google.com/p/goprotobuf/lib/codereview"
list="$list code.google.com/p/goprotobuf/proto"
list="$list code.google.com/p/goprotobuf/proto/testdata"
list="$list code.google.com/p/goprotobuf/protoc-gen-go"
list="$list code.google.com/p/goprotobuf/protoc-gen-go/descriptor"
list="$list code.google.com/p/goprotobuf/protoc-gen-go/generator"
list="$list code.google.com/p/goprotobuf/protoc-gen-go/plugin"
# list="$list code.google.com/p/goprotobuf/protoc-gen-go/testdata"
# list="$list code.google.com/p/goprotobuf/protoc-gen-go/testdata/multi"
list="$list code.google.com/p/goprotobuf/protoc-gen-go/testdata/my_test"
# list="$list code.google.com/p/go-uuid"
# list="$list code.google.com/p/go-uuid/lib"
# list="$list code.google.com/p/go-uuid/lib/codereview"
list="$list code.google.com/p/go-uuid/uuid"
# list="$list github.com"
# list="$list github.com/alphazero"
list="$list github.com/alphazero/Go-Redis"
# list="$list github.com/alphazero/Go-Redis/bench"
# list="$list github.com/alphazero/Go-Redis/compliance"
# list="$list github.com/alphazero/Go-Redis/examples"
# list="$list github.com/alphazero/Go-Redis/test"
# list="$list github.com/alphazero/Go-Redis/test/gen"
# list="$list github.com/garyburd"
# list="$list github.com/garyburd/redigo"
list="$list github.com/garyburd/redigo/redis"
list="$list github.com/garyburd/redigo/redisx"
# list="$list github.com/go-sql-driver"
list="$list github.com/go-sql-driver/mysql"
# list="$list github.com/icattlecoder"
list="$list github.com/icattlecoder/godaemon"
list="$list github.com/icattlecoder/godaemon/example"
# list="$list github.com/vmihailenco"
list="$list github.com/vmihailenco/bufio"
list="$list github.com/vmihailenco/msgpack"
# list="$list gopkg.in"
list="$list gopkg.in/mgo.v2"
list="$list gopkg.in/mgo.v2/bson"
# list="$list gopkg.in/mgo.v2/sasl"
# list="$list gopkg.in/mgo.v2/testdb"
list="$list gopkg.in/mgo.v2/txn"

for name in $list; do
	go install $name
	if [[ $? -ne 0 ]]; then
		echo "=========================================="
		echo "install failed $name"
		echo "=========================================="
		# exit 0
	else
		echo "install ok $name"
	fi
done

# echo "=========================================="
# echo "3rdparty"
# go install 3rdparty
