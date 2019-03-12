package runtime

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/jucardi/go-streams/streams"
	fmtc "github.com/jucardi/go-terminal-colors"
)

type StackFrame struct {
	Line int
	File string
	Fn   string
}

// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
//
// Similar to what `debug.Stack` returns but with a bigger buffer.
func Stack() []byte {
	buf := make([]byte, 2048)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}

// ColoredSimpleStack does the same as `SimpleStack` but uses terminal colors for filename, function and line number
func ColoredSimpleStack(skip ...int) string {
	s := 0
	if len(skip) > 0 {
		s = skip[0]
	}
	return stackMessage(true, s)
}

// SimpleStack is a simpler stack trace that prints the caller lines for each file and line number
// by using the `runtime.Caller`. Unlike `runtime.Stack()` and `debug.Stack()`, it allows to skip
// the N most recent callers, to produce a cleaner stack trace.
func SimpleStack(skip ...int) string {
	s := 0
	if len(skip) > 0 {
		s = skip[0]
	}
	return stackMessage(false, s)
}

// StackFrames returns an array of the stack trace frame data (File, function, line)
func StackFrames(skip ...int) []*StackFrame {
	s := 0
	if len(skip) > 0 {
		s = skip[0]
	}
	frames, _, _ := stackFrames(s)
	return frames
}

// StackLines returns all the lines of code the current context has hit
func StackLines(skip ...int) []string {
	return streams.From(StackFrames(skip...)).
		Filter(func(i interface{}) bool {
			x := i.(*StackFrame)
			return strings.TrimSpace(x.File) != ""
		}).
		Map(func(i interface{}) interface{} {
			x := i.(*StackFrame)
			return fmt.Sprintf("%s:%d", x.File, x.Line)
		}).
		ToArray().([]string)
}

func stackFrames(skip int) (stack []*StackFrame, maxFileLen int, maxFnLen int) {
	rpc := make([]uintptr, 256)
	runtime.Callers(skip, rpc[:])

	frames := runtime.CallersFrames(rpc)
	frame, ok := frames.Next()
	ok = true

	for ; ok; frame, ok = frames.Next() {
		file := getCallerString(frame.File)
		if strings.TrimSpace(frame.File) == "" || strings.HasPrefix(file, "github.com/jucardi/go-titan/runtime") {
			continue
		}
		pkg := strings.Split(file, "/")
		fn := strings.TrimLeft(strings.Replace(frame.Function, strings.Join(pkg[:len(pkg)-1], "/"), "", 1), ".")
		stack = append(stack, &StackFrame{Line: frame.Line, File: file, Fn: fn})
		if maxFileLen < len(file) {
			maxFileLen = len(file)
		}
		if maxFnLen < len(fn) {
			maxFnLen = len(fn)
		}
	}
	stack = stack[skip:]
	return
}

func stackMessage(colored bool, skip int) string {
	stack, maxFileLen, maxFnLen := stackFrames(skip)

	return strings.Join(
		streams.
			From(stack).
			Map(func(i interface{}) interface{} {
				x := i.(*StackFrame)
				return sprintf(colored, "file", "%-"+strconv.Itoa(maxFileLen)+"s", x.File) + " | " +
					sprintf(colored, "fn", "%-"+strconv.Itoa(maxFnLen)+"s", x.Fn) + " | line:" +
					sprintf(colored, "line", "%d", x.Line)
			}).
			ToArray().([]string),
		"\n")
}

func sprintf(colored bool, field, format string, args ...interface{}) string {
	if !colored {
		return fmt.Sprintf(format, args...)
	}
	switch field {
	case "file":
		return fmtc.WithColors(fmtc.Yellow).Sprintf(format, args...)
	case "fn":
		return fmtc.WithColors(fmtc.Cyan).Sprintf(format, args...)
	case "line":
		return fmtc.WithColors(fmtc.White).Sprintf(format, args...)
	}
	return fmt.Sprintf(format, args...)
}
