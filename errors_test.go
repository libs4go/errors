package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallStackError(t *testing.T) {
	println(newCallStackError("newCallStackError", fmt.Errorf("test"), 2).Error())
}

func TestWrapError(t *testing.T) {
	err := fmt.Errorf("test")

	origin := err

	err = newCallStackError("newCallStackError", err, 2)

	require.Equal(t, Cause(err), origin)
	require.Equal(t, Unwrap(err), origin)

	err1 := newCallStackError("newCallStackError", err, 2)

	require.Equal(t, Cause(err1), err)
	require.Equal(t, Unwrap(err1), origin)

	println(err1.Error())
}

func TestPublicApi(t *testing.T) {
	println(Wrap(fmt.Errorf("test"), "test a").Error())
}

type C struct {
}

func (c *C) Hello() {

}

type HelloWold interface {
	Hello()
}

func TestErroCode(t *testing.T) {
	ec := New("test",
		WithCode(-2),
		WithVendor("test"),
		WithAttr("a", "test"), WithAttr("c", &C{}))

	vendor, _ := Vendor(ec)

	require.Equal(t, vendor, "test")

	_, ok := Vendor(fmt.Errorf("test"))

	require.False(t, ok)

	code, _ := Code(ec)

	require.Equal(t, code, -2)

	_, ok = Code(fmt.Errorf("test"))

	require.False(t, ok)

	var a string

	require.True(t, Attr(ec, "a", &a))

	require.Equal(t, a, "test")

	require.False(t, Attr(ec, "b", &a))

	var helloWorld HelloWold

	require.True(t, Attr(ec, "c", &helloWorld))

	require.False(t, Attr(ec, "a", &helloWorld))
}

func TestExample(t *testing.T) {
	err := fmt.Errorf("test")
	err = Wrap(err, "create new stack trace wrapper error")

	// look through callstack
	// StackTrace(err, func(frame runtime.Frame) {
	// 	println(fmt.Sprintf("%s(%s:%d)\n", frame.Function, filepath.Base(frame.File), frame.Line))
	// })

	println(err.Error())
}

func TestAsIs(t *testing.T) {
	err := New("test")
	errWrap := Wrap(err, "create new stack trace wrapper error")

	require.True(t, Is(errWrap, err))

	var target *errorCode

	require.True(t, As(errWrap, &target))

	require.True(t, err == target)
}
