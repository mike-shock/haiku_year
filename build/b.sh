#!/bin/bash
cd ..
DIR=`pwd`
BIN=$DIR/build
OUT=~/D/tmp
DOC=$DIR/docs
LOG=binaries.txt
NAME=haiku_year
NAME0=haiku_year-0

cd $DIR
echo ...Building binary for Android...
cd $BIN
./ba.sh $DIR

echo ...Building binary for GNU/Linux...
cd $BIN
./bl.sh $DIR

cd $DIR
mv ./$NAME.apk $OUT
mv ./$NAME.tar.xz $OUT

cd $OUT
echo ...Extracting Linux binary...
tar -xf $NAME.tar.xz usr/local/bin/$NAME
mv $OUT/usr/local/bin/$NAME $OUT
rm -dr $OUT/usr/

echo ...Repacking executable only...
rm -f $OUT/$NAME.tar.gz
tar czvf $NAME.tar.gz $NAME

echo ...Cross-compiling for MS Windows...
rm -rf $OUT/fyne-cross/
cd $BIN
./bw.sh $DIR
mv $DIR/fyne-cross $OUT
cp $OUT/fyne-cross/bin/windows-amd64/$NAME.exe $OUT

echo ...Packing EXE to ZIP...
cd $OUT
rm -f $OUT/$NAME.exe.zip
zip -9 $NAME.exe.zip $NAME.exe

echo ...Packing initial APK...
rm -f $OUT/$NAME0.apk.zip
zip -9 $NAME0.apk.zip $NAME.apk

echo ...Repacking APK to reduce the size...
rm -rf $OUT/apk/
unzip $NAME.apk -d $OUT/apk
cd $OUT/apk
rm -f $OUT/$NAME.apk
zip -9r $OUT/$NAME.apk ./*
cd $OUT
zip -0 $NAME.apk.zip $NAME.apk

echo ...Cleaning temporary files...
rm -rf $OUT/fyne-cross/
rm -rf $OUT/apk/
rm -f $OUT/$NAME.tar.xz

echo Saving the log...
$DIR/bin/size.sh

echo Copying archived binaries
cp $OUT/$NAME.apk.zip $DIR/deploy
cp $OUT/$NAME.exe.zip $DIR/deploy
cp $OUT/$NAME.tar.gz $DIR/deploy

cd $DIR
