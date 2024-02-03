#!/bin/sh

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

libs=(
    # "golang.org/x/net/html/charset"
    # "golang.org/x/text/transform"
    # "github.com/bmaupin/go-epub"
    # "github.com/nfnt/resize"
    # "github.com/tealeg/xlsx"
    # "github.com/vmihailenco/msgpack"
    # "golang.org/x/sys/windows"
    # "github.com/go-gl/glfw/v3.3/glfw"
    # "github.com/AllenDang/cimgui-go"
    # "github.com/go-gl/gl/v2.1/gl"
)
export go="/c/Program Files/Go/bin/go.exe"
# export go="/usr/local/go/bin/go"
export GOPATH=$DIR/LIB/temp

rm -r -f $DIR/LIB/temp/*
mkdir -p $DIR/LIB/temp



function install_lib {
    for l in "${libs[@]}"
    do
        mkdir -p $DIR/LIB/temp/src
        cd $DIR/LIB/temp/src
        "${go}" mod init "${l//@*/}"
        "${go}" get -u $l
        cp -rf $DIR/LIB/temp/pkg/mod/* $DIR/LIB/src
        rm -r -f $DIR/LIB/src/cache/
        rm -r -f $DIR/LIB/temp/*
    done
}

rm -r -f $DIR/LIB/temp/

install_lib
