package clog

import "log/slog"

type customHandlerOption struct {
	handlerOption  *slog.HandlerOptions
	levelColorMap  map[slog.Level]colorFunc
	isKeysColored  *bool
	isLevelColored *bool
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
