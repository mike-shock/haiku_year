#!/bin/bash
cd ..
DIR=$1
fyne package -os linux -icon $DIR/haiku-year.png
cd $DIR
