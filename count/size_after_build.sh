#!/bin/bash
DIR=~/D/Shock/projects/haiku_year
OUT=~/D/tmp
DOC=$DIR/docs
LOG=binaries.txt

cd $DIR
./ba.sh
./bl.sh

mv ./haiku_year.apk $OUT
mv ./haiku_year.tar.xz $OUT

cd $OUT
tar -xf haiku_year.tar.xz usr/local/bin/haiku_year
mv $OUT/usr/local/bin/haiku_year $OUT
rm -dr $OUT/usr/

ls -l $OUT/haiku* > $OUT/$LOG
cp $OUT/$LOG $DOC

# mv fyne-cross
