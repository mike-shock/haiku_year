set LOG=.\testing.log

go test . -v --count=1 > %LOG%
