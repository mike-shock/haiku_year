#!/bin/bash
DIR=~/D/Shock/projects/haiku_year
OUT=~/D/tmp
DOC=$DIR/docs
LOG=binaries.txt
NAME=haiku_year

ls -l $OUT/haiku* > $OUT/$LOG
cp $OUT/$LOG $DOC
