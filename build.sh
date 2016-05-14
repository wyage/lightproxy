#!/bin/bash

export GOPATH=`pwd`
if [ -d bin ]
then
  echo `date` "bin already exists"
else
  mkdir bin
fi

if [ -z "$GOROOT" ]
then
  echo `date` '$GOROOT not set'
  t=`which go`
  t=`dirname $t`
  export GOROOT=`dirname $t`
fi

echo `date` '$GOROOT='$GOROOT
echo `date` '$GOPATH='$GOPATH

echo `date` start build...
go build -o bin/lightproxy src/main/startup.go

cp src/main/config.json bin/

echo `date` done
echo `date` try edit bin/config.json to suit your need