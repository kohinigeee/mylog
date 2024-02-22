package clog

import (
	"log/slog"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/kohinigeee/mylog/inner/customcoll"
)

type colorFunc func(string, ...any) string
type ClogColor int

const (
	ClogRed ClogColor = iota
	ClogGreen
	ClogBlue
	ClogYellow
	ClogMagenta
	ClogWhite
	ClogHiRed
	ClogHiGreen
	ClogHiBlue
	ClogHiYellow
	ClogHiMagenta
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
	case ClogMagenta:
		return color.MagentaString
	case ClogWhite:
		return color.WhiteString
	case ClogHiRed:
		return color.HiRedString
	case ClogHiGreen:
		return color.HiGreenString
	case ClogHiBlue:
		return color.HiBlueString
	case ClogHiYellow:
		return color.HiYellowString
	case ClogHiMagenta:
		return color.HiMagentaString
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

func splitDirs(path string) []string {
	dirs := strings.Split(path, string(filepath.Separator))
	return dirs
}

func makeLogginSurcepath(fpath string, prefixFoldaName string) string {
	fpath = filepath.FromSlash(fpath)
	dirs := splitDirs(fpath)

	ok, prefixIdx := customcoll.Contains(dirs, prefixFoldaName)

	if !ok {
		return fpath
	}

	prefix := ""
	if prefixIdx > 0 {
		prefix = "..." + string(filepath.Separator)
	}

	dirs = dirs[prefixIdx:]
	logPath := filepath.Join(dirs...)

	return prefix + logPath
}

type clogOrderLevelT struct {
	level uint
}

var (
	orderLevel1 = clogOrderLevelT{level: 1}
	orderLevel2 = clogOrderLevelT{level: 2}
	orderLevel3 = clogOrderLevelT{level: 3}
	orderLevel4 = clogOrderLevelT{level: 4}
)

func removeOrderString(str string) string {
	tokens := strings.Split(str, " ")
	if len(tokens) < 2 {
		return str
	}

	pattern := "!.*!"
	re := regexp.MustCompile(pattern)
	if re.MatchString(tokens[0]) {
		return tokens[1]
	}

	return str
}

func OrderLevel1() clogOrderLevelT {
	return orderLevel1
}
func OrderLevel2() clogOrderLevelT {
	return orderLevel2
}
func OrderLevel3() clogOrderLevelT {
	return orderLevel3
}
func OrderLevel4() clogOrderLevelT {
	return orderLevel4
}

func OrderString(keyName string, order clogOrderLevelT) string {
	orderStr := strconv.Itoa(int(order.level + 1))

	return "!" + orderStr + "! " + keyName
}
