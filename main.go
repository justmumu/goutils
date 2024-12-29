package main

import "github.com/justmumu/goutils/logutil"

func main() {
	conf := logutil.LoggerConfig{
		ConsoleEnabled: true,
		ConsoleLevel:   logutil.InfoLevel,
		ConsoleJson:    false,

		FileEnabled: true,
		FileLevel:   logutil.DebugLevel,
		FileJson:    false,

		LogDirectory: "/Users/rasitaydin/Desktop/DEVELOPMENTS/goutils",
		Filename:     "test.log",
		MaxSize:      1,
		MaxBackup:    3,
		MaxAge:       1,
	}

	l := logutil.NewLogger(conf)

	l.Info("TEST1", "TEST2")
	l.Infof("TEST3 %s", "test")
	l.Infoln("TEST4", "TEST5")
	l.Infow("TEST6", "key1", "value1", "key2", "value2")

	l.Debug("TEST1", "TEST2")
	l.Debugf("TEST3 %s", "test")
	l.Debugln("TEST4", "TEST5")
	l.Debugw("TEST6", "key1", "value1", "key2", "value2")
}
