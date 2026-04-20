#!/bin/bash
PRG=../constants.go
DATA=../haiku/year
LOG=../docs/count.txt
TOTAL=367
NOW=`ls -R -1 $DATA/ | grep '[0-9][0-9]-[0-9][0-9].txt' | wc -l`
((LEFT=$TOTAL - $NOW))
TODAY=`date +"%Y-%m-%d"`
VARIANTS=`find $DATA -type f -name "*.txt" | wc -l`
((V=$VARIANTS - 1))

echo                      >> $LOG
echo $TODAY:              >> $LOG
echo -e '\t'Days: $TOTAL  >> $LOG
echo -e '\t'Now: $NOW     >> $LOG
echo -e '\t'Variants: $V  >> $LOG
echo -e '\t'Left: $LEFT   >> $LOG

echo "package main"          > $PRG
echo "const ("              >> $PRG
echo "	HAIKU_TOTAL = $V"   >> $PRG
echo "	HAIKU_LEFT = $LEFT" >> $PRG
echo ")"                    >> $PRG
