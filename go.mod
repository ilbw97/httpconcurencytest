module httpconcurencytest

go 1.17

replace github.com/ilbw97/debuglog => ../debuglog

require (
	github.com/google/brotli/go/cbrotli v0.0.0-20220512075048-9801a2c5d6c6
	github.com/ilbw97/debuglog v0.0.0-20220720072703-f8d4aca3dd9e
	github.com/sirupsen/logrus v1.9.0
)

require (
	golang.org/x/sys v0.0.0-20220808155132-1c4a2a72c664 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)
