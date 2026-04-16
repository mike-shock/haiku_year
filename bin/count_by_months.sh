#!/bin/bash

BIN=`pwd`
DATA=$BIN/../haiku/year
LOG=$BIN/../docs/count_by_months.txt

TODAY=`date +"%Y-%m-%d"`
echo                      >> $LOG
echo $TODAY:              >> $LOG

cd $DATA
echo -e '\t'"Month"'\t'Final'\t'Variants >> $LOG
for dir in */; do
  F=`ls -1 $dir/ | grep '[0-9][0-9]-[0-9][0-9].txt' | wc -l`
  A=`ls -1 $dir/*.txt | wc -l`
  echo -e '\t'$dir'\t'$F'\t'$A >> $LOG
done

