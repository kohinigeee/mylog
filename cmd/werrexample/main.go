package main

import (
	"fmt"
	"os"

	"github.com/kohinigeee/mylog/werr"
)

var (
	w *werr.Wrapper
)

func testFunc() {
	err := w.Errf(0, "test error")
	fmt.Printf("%v\n", err)

	err2 := w.WrapErrf(0, err, "wrap error")
	fmt.Printf("%v\n", err2)

	err3 := w.WrapErrf(0, err2, "wrap error2")
	fmt.Printf("%v\n", err3)
}

func main() {
	cwd, _ := os.Getwd()
	fmt.Printf("cwd: %s\n", cwd)
	w = werr.NewWrapper(
		werr.WithPrefixDir(cwd),
		werr.WithPrefixModuleName("main"),
	)

	testFunc()
}
