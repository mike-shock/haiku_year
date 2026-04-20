#!/bin/bash

BIN=`pwd`
DATA=$BIN/../haiku/year
LOG=$BIN/../docs/count_by_months.txt

declare -A months
months["00/"]=1
months["01/"]=31
months["02/"]=29
months["03/"]=31
months["04/"]=30
months["05/"]=31
months["06/"]=30
months["07/"]=31
months["08/"]=31
months["09/"]=30
months["10/"]=31
months["11/"]=30
months["12/"]=31

TODAY=`date +"%Y-%m-%d"`
echo                      >> $LOG
echo $TODAY:              >> $LOG

cd $DATA
echo -e '\t'"Month"'\t'Final'\t'Variants'\t'Days'\t'Left >> $LOG
for dir in */; do
  F=`ls -1 $dir/ | grep '[0-9][0-9]-[0-9][0-9].txt' | wc -l`
  A=`ls -1 $dir/*.txt | wc -l`
  M=${months[$dir]}
  ((L=$M - $F))
  echo -e '\t'$dir'\t'$F'\t'$A'\t\t'$M'\t'$L >> $LOG
done

