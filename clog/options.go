package clog

import "log/slog"

type AddSourceOption struct {
	PrefixFoldaName string
	AddSource       bool
}

type customHandlerOption struct {
	handlerOption  *slog.HandlerOptions
	levelColorMap  map[slog.Level]colorFunc
	isShow         *bool
	isKeysColored  *bool
	isLevelColored *bool
	levelFilter    map[slog.Level]bool
	addSourceOpt   AddSourceOption
}

type customHandlerOptionFunc func(options *customHandlerOption) error

func WithHandlerOption(option *slog.HandlerOptions) customHandlerOptionFunc {
	return func(options *customHandlerOption) error {
		options.handlerOption = option
		return nil
	}
}

func WithSetLevelColor(level slog.Level, color ClogColor) customHandlerOptionFunc {
	return func(options *customHandlerOption) error {
		options.levelColorMap[level] = getColorFunc(color)
		return nil
	}
}

func WithKeyColored(isKeysColored bool) customHandlerOptionFunc {
	return func(options *customHandlerOption) error {
		options.isKeysColored = &isKeysColored
		return nil
	}
}

func WithLevelColored(isLevelColored bool) customHandlerOptionFunc {
	return func(options *customHandlerOption) error {
		options.isLevelColored = &isLevelColored
		return nil
	}
}

func WithShow(isShow bool) customHandlerOptionFunc {
	return func(option *customHandlerOption) error {
		option.isShow = &isShow
		return nil
	}
}

func WithAddFilter(level slog.Level) customHandlerOptionFunc {
	return func(option *customHandlerOption) error {
		option.levelFilter[level] = true
		return nil
	}
}

func WithAddSourceOption(opt AddSourceOption) customHandlerOptionFunc {
	return func(option *customHandlerOption) error {
		option.addSourceOpt = opt
		return nil
	}
}
