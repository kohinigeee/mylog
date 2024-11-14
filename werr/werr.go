package werr

import "fmt"

type werr struct {
	file     string
	line     int
	funcName string
	msg      string
	err      error
	depth    int
}

func newWerr(file string, line int, funcName string, fmtstr string, args ...interface{}) werr {

	msg := fmt.Sprintf(fmtstr, args...)
	depth := 0

	// errmsg := fmt.Sprintf("{ [msg] %s | [run] %s %s(L%d) }", msgerr, file, funcName, line)
	err := fmt.Errorf("{ msg[%d] %s | run[%d] %s %s(L%d) }", depth, msg, depth, file, funcName, line)

	return werr{
		file:     file,
		line:     line,
		funcName: funcName,
		msg:      msg,
		err:      err,
		depth:    depth,
	}
}

func newWerrByWrap(file string, line int, funcName string, originErr werr, fmtstr string, args ...interface{}) werr {
	msg := fmt.Sprintf(fmtstr, args...)
	depth := originErr.depth + 1

	err := fmt.Errorf("{ msg[%d] %s | run[%d] %s %s(L%d) | trace[%d]-> %w }", depth, msg, depth, file, funcName, line, depth, originErr.err)

	return werr{
		file:     file,
		line:     line,
		funcName: funcName,
		msg:      msg,
		err:      err,
		depth:    depth,
	}
}

func (w werr) File() string {
	return w.file
}

func (w werr) Line() int {
	return w.line
}

func (w werr) FuncName() string {
	return w.funcName
}

func (w werr) Msg() string {
	return w.msg
}

func (w werr) Error() string {
	return w.err.Error()
}
