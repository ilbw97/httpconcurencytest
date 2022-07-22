module httpconcurencytest

go 1.18

require (
	github.com/ilbw97/debuglog v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.0
)

require (
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace github.com/ilbw97/debuglog => ../debuglog
