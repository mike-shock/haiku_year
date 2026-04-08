#!/bin/bash
DIR=~/D/Shock/projects/haiku_year
OUT=~/D/tmp
DOC=$DIR/docs
LOG=binaries.txt
NAME=haiku_year

cd $DIR
./ba.sh
./bl.sh
./bw.sh

mv ./$NAME.apk $OUT
mv ./$NAME.tar.xz $OUT

cd $OUT
tar -xf $NAME.tar.xz usr/local/bin/$NAME
mv $OUT/usr/local/bin/$NAME $OUT
rm -dr $OUT/usr/

rm -rf $OUT/fyne-cross/
mv $DIR/fyne-cross $OUT
cp $OUT/fyne-cross/bin/windows-amd64/$NAME.exe $OUT

cd $OUT
rm -f $OUT/$NAME.exe.zip
zip -9 $NAME.exe.zip $NAME.exe

rm -f $OUT/$NAME.apk.zip
zip -9 $NAME.apk.zip $NAME.apk

rm -f $OUT/$NAME.apk.7z
7z a $NAME.apk.7z $NAME.apk

ls -l $OUT/haiku* > $OUT/$LOG
cp $OUT/$LOG $DOC
