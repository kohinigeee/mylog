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

	myh, err := clog.NewCustomTextHandler(os.Stdout,
		clog.WithHandlerOption(&slog.HandlerOptions{
			Level: logLevel,
		}),
		clog.WithAddSourceOption(clog.AddSourceOption{
			PrefixFoldaName: "mylog",
			AddSource:       true,
			CallerDepth:     0,
		}),
	)

	if err != nil {
		panic(err)
	}

	myLogger := slog.New(myh)

	var s upperCase = "it is upper case string"

	myLogger.Debug("Hello World",
		"upperCase", s,
	)

	myLogger.Info("Hello World")
	myLogger.Warn("Hello World", slog.String("key", "value"), slog.String("123", "value"))
	myLogger.Error("Hello World", slog.String("key", "value"))
}
