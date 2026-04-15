#!/bin/bash
cd ..
DIR=$1
export ANDROID_NDK_HOME=$HOME/Android/Sdk/ndk/29.0.14206865
fyne package -os android -app-id com.shokhirev.mike.haiku_year -icon $DIR/haiku-year.png --app-version 0.3.30 --release
cd $DIR

