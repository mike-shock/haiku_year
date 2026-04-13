#!/bin/bash

BIN=`pwd`
DATA=$BIN/../haiku/year
LOG=$BIN/../docs/count_by_months.txt

TODAY=`date +"%Y-%m-%d"`
echo                      >> $LOG
echo $TODAY:              >> $LOG

cd $DATA
for dir in */; do
  D=`ls -1 $dir/ | grep '[0-9][0-9]-[0-9][0-9].txt' | wc -l`
  echo -e '\t'"$dir" $D >> $LOG
done

