package errors

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type callStack struct {
	stackPC []uintptr
}

func newCallStack(skip int) *callStack {
	pcs := make([]uintptr, 32)

	count := runtime.Callers(skip, pcs)

	return &callStack{
		stackPC: pcs[:count],
	}
}

func (cs *callStack) String() string {

	var buff bytes.Buffer

	cs.WalkCallFrames(func(frame runtime.Frame) {
		buff.WriteString(fmt.Sprintf("    at %s(%s:%d)\n", frame.Function, filepath.Base(frame.File), frame.Line))
	})

	return buff.String()
}

func (cs *callStack) WalkCallFrames(callback func(runtime.Frame)) {
	frames := runtime.CallersFrames(cs.stackPC)

	for {
		frame, more := frames.Next()

		if index := strings.Index(frame.File, "src"); index != -1 {
			// trim GOPATH or GOROOT prifix
			frame.File = string(frame.File[index+4:])
		}

		callback(frame)

		if !more {
			break
		}
	}
}
