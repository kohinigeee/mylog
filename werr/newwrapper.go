package werr

import (
	"fmt"
	"strings"
)

type werrNewOptionalFunc func(w *Wrapper) error

func NewWrapper(options ...werrNewOptionalFunc) *Wrapper {
	w := &Wrapper{
		prefixDir:     "",
		pathSeparater: Slash,
	}

	for _, opt := range options {
		err := opt(w)
		if err != nil {
			panic(fmt.Errorf("[NewWrapper] error setting option : %w", err))
		}
	}
	return w
}

// Set prefix directory(if you set prefix directory, showed path will change to relative path from prefix directory)
func WithPrefixDir(prefixDir string) werrNewOptionalFunc {
	return func(w *Wrapper) error {
		w.prefixDir = w.toInPath(prefixDir)
		if !strings.HasSuffix(string(w.prefixDir), "/") {
			w.prefixDir += "/"
		}
		return nil
	}
}

// Change separeter type of path(BackSlash or Slash)
func WithPathSeparater(pathSeparater pathSeparaterType) werrNewOptionalFunc {
	return func(w *Wrapper) error {
		w.pathSeparater = pathSeparater
		w.prefixDir = w.toInPath(string(w.prefixDir))
		return nil
	}
}
