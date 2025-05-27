package evil

import "runtime"

type fnInfo struct {
	entry    uintptr
	name     string
	pkg      string
	line     int
	filename string
}

type callerInfo struct {
	pc       uintptr
	line     int
	pkg      string
	filename string

	fn *fnInfo
}

func getCallerInfo(pc int) *callerInfo {
	var ok bool

	ci := new(callerInfo)
	ci.pc, ci.filename, ci.line, ok = runtime.Caller(pc + 1)
	if !ok {
		panic("could not get calling function's PC")
	}

	callerFn := runtime.FuncForPC(ci.pc)

	ci.fn = &fnInfo{
		name:  callerFn.Name(),
		entry: callerFn.Entry(),
	}
	ci.fn.filename, ci.fn.line = callerFn.FileLine(ci.pc)

	return ci
}
