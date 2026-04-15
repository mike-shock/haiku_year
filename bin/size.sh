#!/bin/bash
DIR=~/D/Shock/projects/haiku_year
OUT=~/D/tmp
DOC=$DIR/docs
LOG=binaries.txt

ls -l $OUT/haiku* > $OUT/$LOG
cp $OUT/$LOG $DOC
