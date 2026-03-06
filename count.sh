#!/bin/bash
ls -R -1 haiku/ | grep '[0-9][0-9]-[0-9][0-9].txt' | wc -l
