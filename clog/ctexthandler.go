package clog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/kohinigeee/mylog/inner/customcoll"
)

type customTextHandler struct {
	slog.Handler
	w              io.Writer
	levelColorMap  map[slog.Level]colorFunc
	isKeysColored  bool
	isLevelColored bool
	isShow         bool
	filterLevel    map[slog.Level]bool
	addSourceOpt   AddSourceOption
}

func NewCustomTextHandler(w io.Writer, opts ...customHandlerOptionFunc) (*customTextHandler, error) {
	option := &customHandlerOption{
		handlerOption:  nil,
		levelColorMap:  defaultColorMap(),
		isKeysColored:  nil,
		isLevelColored: nil,
		isShow:         nil,
		levelFilter:    make(map[slog.Level]bool),
		addSourceOpt: AddSourceOption{
			PrefixFoldaName: "",
			AddSource:       false,
		},
	}

	for _, opt := range opts {
		err := opt(option)
		if err != nil {
			return nil, fmt.Errorf("error in customTextHandlerOptionFunc: %w", err)
		}
	}

	if option.handlerOption == nil {
		option.handlerOption = &slog.HandlerOptions{}
	}

	if option.isShow == nil {
		initialShow := true
		option.isShow = &initialShow
	}

	if option.isKeysColored == nil {
		initialColored := true
		option.isKeysColored = &initialColored
	}

	if option.isLevelColored == nil {
		initialColored := true
		option.isLevelColored = &initialColored
	}

	return &customTextHandler{
		Handler:        slog.NewTextHandler(w, option.handlerOption),
		w:              w,
		levelColorMap:  option.levelColorMap,
		isShow:         *option.isShow,
		isKeysColored:  *option.isKeysColored,
		isLevelColored: *option.isLevelColored,
		filterLevel:    option.levelFilter,
		addSourceOpt:   option.addSourceOpt,
	}, nil
}

func (h *customTextHandler) filedsNameToColorText(name string, nestLevel int) string {
	if !h.isKeysColored {
		return name
	}

	nestColorMap := map[int]colorFunc{
		0: color.BlueString,
		1: color.GreenString,
		2: color.YellowString,
	}
	clorFunc := nestColorMap[nestLevel%len(nestColorMap)]
	return clorFunc(name)
}

func (h *customTextHandler) filedsToText(fields map[string]any, prefixs []string) []string {
	ans := make([]string, 0)

	keys := customcoll.MapKeys[string](fields)
	sort.Strings(keys)

	for _, key := range keys {
		value := fields[key]

		keyName := removeOrderString(key)
		if _, ok := value.(map[string]any); ok {
			keyName = h.filedsNameToColorText(keyName, len(prefixs))
			prefixs = append(prefixs, keyName)

			ans = append(ans, h.filedsToText(value.(map[string]any), prefixs)...)
			prefixs = prefixs[:len(prefixs)-1]
		} else {
			keyName := h.filedsNameToColorText(keyName, len(prefixs))
			prefixs = append(prefixs, keyName)
			prefixName := strings.Join(prefixs, ".")

			text := fmt.Sprintf("%v:%v", prefixName, value)

			ans = append(ans, text)
			prefixs = prefixs[:len(prefixs)-1]
		}
	}

	return ans
}

func (h *customTextHandler) textWithLevel(str string, level slog.Level) string {
	if !h.isLevelColored {
		return str
	}

	if colorFunc, ok := h.levelColorMap[level]; ok {
		return colorFunc(str)
	}
	return str
}

func (h *customTextHandler) Enabled(c context.Context, level slog.Level) bool {
	if !h.isShow {
		return false
	}

	if len(h.filterLevel) > 0 {
		_, ok := h.filterLevel[level]
		return ok
	}

	return h.Handler.Enabled(c, level)
}

func (h *customTextHandler) Handle(_ context.Context, r slog.Record) error {

	timeStr := r.Time.Format("2006-01-02 15:04:05")
	levelStr := h.textWithLevel(r.Level.String(), r.Level)
	msg := h.textWithLevel(r.Message, r.Level)

	fieldsMap := make(map[string]any)
	r.Attrs(func(a slog.Attr) bool {
		addFields(fieldsMap, a)
		return true
	})

	if h.addSourceOpt.AddSource {
		const skipLevel = 3
		_, file, line, _ := runtime.Caller(skipLevel)
		file = makeLogginSurcepath(file, h.addSourceOpt.PrefixFoldaName)
		logStr := fmt.Sprintf("%v:%v", file, line)
		keyName := OrderString("logCode", OrderLevel1())
		addFields(fieldsMap, slog.String(keyName, logStr))
	}

	fieldsTexts := h.filedsToText(fieldsMap, []string{})

	logText := ""
	if h.isLevelColored {
		logText += fmt.Sprintf("\n%v [%-15s] %v", timeStr, levelStr, msg)
	} else {
		logText += fmt.Sprintf("\n%v [%-6s] %v", timeStr, levelStr, msg)
	}
	if len(fieldsTexts) > 0 {
		logText += fmt.Sprintf("\n  %v", strings.Join(fieldsTexts, "  "))
	}
	logText += "\n"

	h.w.Write([]byte(logText))
	return nil
}
