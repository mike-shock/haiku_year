#!/bin/bash
TOTAL=366
NOW=`ls -R -1 haiku/ | grep '[0-9][0-9]-[0-9][0-9].txt' | wc -l`
echo Now: $NOW
((LEFT=$TOTAL - $NOW))
echo Left: $LEFT
