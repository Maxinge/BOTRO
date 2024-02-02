#!/bin/sh
# install a c++ compiler for GO_CLIB

if [[ $1 == "" ]]; then exit; fi

echo "building target [_"$1"]"

NAME=$1
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export go="/c/Program Files/Go/bin/go.exe"
export GO111MODULE=off
export PATH=$PATH:"/c/Program Files/Go/bin/";
export PATH=$PATH:"/c/Program Files/mingw-w64/x86_64-8.1.0-posix-seh-rt_v6-rev0/mingw64/bin";

export GOPATH=$DIR/LIB/:$DIR/_$NAME
export GOOS=windows
# export GOARCH=386

# echo $GOPATH

if [[ $NAME == "GO_CLIB" ]]; then
    "${go}" build -o _$NAME/$NAME.dll -buildmode=c-shared $DIR/_$NAME/src/main.go
    exit
fi

rm _$NAME/$NAME.exe

source "${DIR}/_${NAME}/src/_conf.sh"

SRCS=""
for X in "${SRC_FILES[@]}"
    do
        SRCS=$SRCS" "${X}
    done

"${go}" build -p 4 -o _$NAME/$NAME.exe ${SRCS}

_$NAME/$NAME.exe

# rm _$NAME/$NAME".exe"

# go build -o _$NAME/$NAME".exe" $DIR/_$NAME/src/main.go
