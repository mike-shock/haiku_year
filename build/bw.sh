#!/bin/bash
cd ..
DIR=$1
GOTOOLCHAIN=auto fyne-cross windows -app-id com.shokhirev.mike.haiku_year -icon haiku-year.png
cd $DIR
