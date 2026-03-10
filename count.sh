#!/bin/bash
TOTAL=366
NOW=`ls -R -1 haiku/ | grep '[0-9][0-9]-[0-9][0-9].txt' | wc -l`
((LEFT=$TOTAL - $NOW))
TODAY=`date +"%Y-%m-%d"`

echo             >> count.txt
echo $TODAY:            >> count.txt
echo -e '\t'Now: $NOW   >> count.txt
echo -e '\t'Left: $LEFT >> count.txt

