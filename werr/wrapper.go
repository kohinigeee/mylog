package werr

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

type pathSeparaterType int

const (
	Slash pathSeparaterType = iota
	BackSlash
)

func (p pathSeparaterType) String() string {
	switch p {
	case Slash:
		return "/"
	case BackSlash:
		return "\\"
	default:
		panic(fmt.Sprintf("unknown pathSeparaterType: %d", p))
	}
}

func toPath(path string, sep pathSeparaterType) string {
	switch sep {
	case Slash:
		return filepath.ToSlash(path)
	case BackSlash:
		return filepath.FromSlash(path)
	default:
		log.Fatalf("unknown pathSeparaterType: %d", sep)
	}
	return ""
}

//----------------------------------------------

type slashPathT string

type Wrapper struct {
	prefixDir     slashPathT
	pathSeparater pathSeparaterType
}

func (w *Wrapper) toOutPath(path slashPathT) string {
	return toPath(string(path), w.pathSeparater)
}

func (w *Wrapper) toInPath(path string) slashPathT {
	return slashPathT(toPath(path, Slash))
}

func (w *Wrapper) toRel(file slashPathT) slashPathT {
	if !strings.HasPrefix(string(file), string(w.prefixDir)) {
		return file
	}

	fstr := strings.TrimPrefix(string(file), string(w.prefixDir))
	return slashPathT("." + Slash.String() + fstr)
}

func (w *Wrapper) PrefixDir() string {
	return string(w.prefixDir)
}

// trace : 0 is current function, 1 is caller function
func (w *Wrapper) Errf(trace int, format string, args ...interface{}) Werr {
	const pcDepth = 1
	pc, fileStr, line, ok := runtime.Caller(trace + pcDepth)

	file := w.toInPath(fileStr)
	if !ok {
		msg := "[Wrapper .Errf] runtime.Caller failed"
		panic(msg)
	}

	funcName := runtime.FuncForPC(pc).Name()
	file = w.toRel(file)

	return newWerr(w.toOutPath(file), line, funcName, format, args...)
}

// trace : 0 is current function, 1 is caller function
func (w *Wrapper) WrapErrf(trace int, err error, format string, args ...interface{}) Werr {
	const pcDepth = 1
	pc, fileStr, line, ok := runtime.Caller(trace + pcDepth)

	file := w.toInPath(fileStr)
	if !ok {
		msg := "[Wrapper .Errf] runtime.Caller failed"
		panic(msg)
	}

	funcName := runtime.FuncForPC(pc).Name()
	file = w.toRel(file)

	return newWerrByWrap(w.toOutPath(file), line, funcName, err, format, args...)
}
