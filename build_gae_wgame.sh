#!/bin/bash

. dev_setup.sh
#go build -ldflags "-X main.version=`date -u +.%Y%m%d.%H%M%S`" gae_wgame

BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')
VERSIONFILE=src/gae_wgame/version.go

rm -f $VERSIONFILE
echo "package main" > $VERSIONFILE
echo "const (" >> $VERSIONFILE
echo "  //BuildVersion version string " >> $VERSIONFILE
echo "  BuildVersion = \"1.0\"" >> $VERSIONFILE
echo "  //BuildTime build date and time " >> $VERSIONFILE
echo "  BuildTime = \"$BUILD_TIME\"" >> $VERSIONFILE
echo ")" >> $VERSIONFILE


go build -o bin/gae_wgame gae_wgame
