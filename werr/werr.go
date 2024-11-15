package werr

import "fmt"

type Werr struct {
	file     string
	line     int
	funcName string
	msg      string
	err      error
	depth    int
}

func newWerr(file string, line int, funcName string, fmtstr string, args ...interface{}) Werr {

	msg := fmt.Sprintf(fmtstr, args...)
	depth := 0

	// errmsg := fmt.Sprintf("{ [msg] %s | [run] %s %s(L%d) }", msgerr, file, funcName, line)
	err := fmt.Errorf("{ msg[%d] %s | run[%d] %s %s(L%d) }", depth, msg, depth, file, funcName, line)

	return Werr{
		file:     file,
		line:     line,
		funcName: funcName,
		msg:      msg,
		err:      err,
		depth:    depth,
	}
}

func newWerrByWrap(file string, line int, funcName string, originErr error, fmtstr string, args ...interface{}) Werr {
	msg := fmt.Sprintf(fmtstr, args...)

	switch originWerr := originErr.(type) {
	case Werr:
		depth := originWerr.depth + 1

		err := fmt.Errorf("{ msg[%d] %s | run[%d] %s %s(L%d) | trace[%d]-> %w }", depth, msg, depth, file, funcName, line, depth, originWerr.err)

		return Werr{
			file:     file,
			line:     line,
			funcName: funcName,
			msg:      msg,
			err:      err,
			depth:    depth,
		}
	default:
		depth := 0
		err := fmt.Errorf("{ msg[%d] %s | run[%d] %s %s(L%d) }", depth, msg, depth, file, funcName, line)

		return Werr{
			file:     file,
			line:     line,
			funcName: funcName,
			msg:      msg,
			err:      err,
			depth:    depth,
		}
	}
}

func (w Werr) File() string {
	return w.file
}

func (w Werr) Line() int {
	return w.line
}

func (w Werr) FuncName() string {
	return w.funcName
}

func (w Werr) Msg() string {
	return w.msg
}

func (w Werr) Error() string {
	return w.err.Error()
}
