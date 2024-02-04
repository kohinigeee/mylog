package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/kohinigeee/mylog/clog"
)

type upperCase string

func (x upperCase) LogValue() slog.Value {
	return slog.StringValue(strings.ToUpper(string(x)))
}

func main() {
	var logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)

	myh, err := clog.NewCustomJsonHandler(os.Stdout,
		clog.WithKeyColored(true),
		clog.WithLevelColored(true),
		clog.WithAddFilter(slog.LevelDebug),
	)

	if err != nil {
		panic(err)
	}

	myLogger := slog.New(myh)

	var s upperCase = "it is upper case string"

	myLogger.Debug("Hello World",
		"upperCase", s,
	)
	myLogger.Info("Hello World", slog.String("key", "value"))
	myLogger.Warn("Hello World")
	myLogger.Error("Hello World", slog.String("key", "value"))
}
