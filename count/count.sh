#!/bin/bash
DATA=../haiku/year
LOG=../docs/count.txt
TOTAL=367
NOW=`ls -R -1 $DATA/ | grep '[0-9][0-9]-[0-9][0-9].txt' | wc -l`
((LEFT=$TOTAL - $NOW))
TODAY=`date +"%Y-%m-%d"`

echo                      >> $LOG
echo $TODAY:              >> $LOG
echo -e '\t'Now: $NOW     >> $LOG
echo -e '\t'Left: $LEFT   >> $LOG
echo -e '\t'Total: $TOTAL >> $LOG

