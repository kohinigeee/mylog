package clog

import (
	"log/slog"

	"github.com/fatih/color"
)

type colorFunc func(string, ...any) string
type ClogColor int

const (
	ClogRed ClogColor = iota
	ClogGreen
	ClogBlue
	ClogYellow
)

func getColorFunc(c ClogColor) colorFunc {
	switch c {
	case ClogRed:
		return color.RedString
	case ClogGreen:
		return color.GreenString
	case ClogBlue:
		return color.BlueString
	case ClogYellow:
		return color.YellowString
	}
	return color.WhiteString
}

func defaultColorMap() map[slog.Level]colorFunc {
	return map[slog.Level]colorFunc{
		slog.LevelDebug: color.BlueString,
		slog.LevelInfo:  color.GreenString,
		slog.LevelWarn:  color.YellowString,
		slog.LevelError: color.RedString,
	}
}

func addFields(fields map[string]any, a slog.Attr) {
	value := a.Value.Any()

	if _, ok := value.([]slog.Attr); !ok {
		key := a.Key

		if _, ok := value.(slog.LogValuer); ok {
			valueStr := value.(slog.LogValuer).LogValue().String()
			fields[key] = valueStr
		} else {
			fields[key] = value
		}
		return
	}

	attrs := value.([]slog.Attr)

	innerFields := make(map[string]any)
	for _, attr := range attrs {
		addFields(innerFields, attr)
	}
	key := a.Key
	fields[key] = innerFields
}

func coloredStrWithLevel(level slog.Level, str string) string {
	levelColorMap := map[slog.Level]colorFunc{
		slog.LevelDebug: color.BlueString,
		slog.LevelInfo:  color.GreenString,
		slog.LevelWarn:  color.YellowString,
		slog.LevelError: color.RedString,
	}

	if colorFunc, ok := levelColorMap[level]; ok {
		return colorFunc(str)
	}
	return str
}

func levelStrWithColor(level slog.Level) string {
	str := level.String()
	return coloredStrWithLevel(level, str)
}
