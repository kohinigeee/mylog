package clog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/kohinigeee/mylog/inner/customcoll"
)

type customJsonHandler struct {
	slog.Handler
	w              io.Writer
	levelColorMap  map[slog.Level]colorFunc
	isKeysColored  bool
	isLevelColored bool
	isShow         bool
	filterLevel    map[slog.Level]bool
}

func NewCustomJsonHandler(w io.Writer, opts ...customHandlerOptionFunc) (*customJsonHandler, error) {

	option := &customHandlerOption{
		handlerOption:  nil,
		levelColorMap:  defaultColorMap(),
		isKeysColored:  nil,
		isLevelColored: nil,
		isShow:         nil,
		levelFilter:    make(map[slog.Level]bool),
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

	return &customJsonHandler{
		Handler:        slog.NewJSONHandler(w, option.handlerOption),
		w:              w,
		levelColorMap:  option.levelColorMap,
		isShow:         *option.isShow,
		isKeysColored:  *option.isKeysColored,
		isLevelColored: *option.isLevelColored,
		filterLevel:    option.levelFilter,
	}, nil
}

func (h *customJsonHandler) myJsonMarshalKeyWithColor(key string, nlevel int) string {
	if !h.isKeysColored {
		return key
	}

	colors := []colorFunc{
		color.BlueString,
		color.GreenString,
		color.YellowString,
	}

	colorIndex := nlevel % len(colors)
	return colors[colorIndex](key)
}

func (h *customJsonHandler) textWithLevel(str string, level slog.Level) string {
	if !h.isLevelColored {
		return str
	}

	if colorFunc, ok := h.levelColorMap[level]; ok {
		return colorFunc(str)
	}
	return str
}

func (h *customJsonHandler) myJsonMarshalNest(fields map[string]any, nlevel int, prefix string) (string, error) {
	prefixIndent := strings.Repeat(prefix, (nlevel+1)*2)

	keys := customcoll.MapKeys(fields)
	//最初のネストの場合はtime, level, msgを先頭に持ってくる
	if nlevel == 0 {
		targetKeys := []string{"time", "level", "msg"}

		keys = customcoll.Filter[string](keys, func(key string) bool {
			if ok, _ := customcoll.Contains(targetKeys, key); ok {
				return false
			}
			return true
		})
		sort.Strings(keys)
		keys = append(targetKeys, keys...)
	} else {
		sort.Strings(keys)
	}

	ans := ""
	ans += "{\n"
	for _, key := range keys {
		value := fields[key]
		keyStr := h.myJsonMarshalKeyWithColor(key, nlevel)
		nameStr := fmt.Sprintf("\"%v\"", keyStr)
		innerStr := ""
		err := error(nil)
		if innerJson, ok := value.(map[string]any); ok {
			innerStr, err = h.myJsonMarshalNest(innerJson, nlevel+1, prefix)
			if err != nil {
				continue
			}
		} else {
			switch value.(type) {
			case string:
				innerStr = fmt.Sprintf("\"%v\"", value)
			case int:
				innerStr = fmt.Sprintf("%v", value)
			default:
				innerStr = fmt.Sprintf("%v", value)
			}
		}
		ans += fmt.Sprintf("%v%v: %v,\n", prefixIndent, nameStr, innerStr)
	}
	ans = strings.TrimSuffix(ans, ",\n")
	nestPrefixStr := strings.Repeat(prefix, nlevel*2)
	ans += fmt.Sprintf("\n%v}", nestPrefixStr)
	return ans, nil
}

func (h *customJsonHandler) Enabled(c context.Context, level slog.Level) bool {
	if !h.isShow {
		return false
	}

	if len(h.filterLevel) > 0 {
		_, ok := h.filterLevel[level]
		return ok
	}

	return h.Handler.Enabled(c, level)
}

func (h *customJsonHandler) Handle(_ context.Context, r slog.Record) error {

	msgStr := h.textWithLevel(r.Message, r.Level)
	timeStr := r.Time.Format("2006-01-02 15:04:05")
	levelStr := h.textWithLevel(r.Level.String(), r.Level)
	fields := make(map[string]any)

	fields["time"] = timeStr
	fields["level"] = levelStr
	fields["msg"] = msgStr

	r.Attrs(func(a slog.Attr) bool {
		addFields(fields, a)
		return true
	})

	jsonStr, err := h.myJsonMarshalNest(fields, 0, " ")
	jsonData := []byte(jsonStr)

	if err != nil {
		fmt.Println("Error marshalling json")
		return err
	}

	h.w.Write([]byte(jsonData))
	return err
}
