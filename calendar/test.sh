#!/bin/bash
LOG=./testing.log

go test . -v --count=1 > $LOG
tail -n 1 $LOG
