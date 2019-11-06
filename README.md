# errors

The substitute for the golang [pkg/errors](https://golang.org/pkg/errors/), support print call stack and more ..

## Wrap exists error object

errors library support wrap exists error object to support stack trace

```go
func TestExample(t *testing.T) {
	err := fmt.Errorf("test")
	err = Wrap(err, "create new stack trace wrapper error")

	// look through callstack
	// StackTrace(err, func(frame runtime.Frame) {
	// 	println(fmt.Sprintf("%s(%s:%d)\n", frame.Function, filepath.Base(frame.File), frame.Line))
	// })

	println(err.Error())
}
```

above codes will print like this:

```txt
error: create new stack trace wrapper error
    at github.com/BlockchainMiddleware/errors.newCallStackError(errors.go:21)
    at github.com/BlockchainMiddleware/errors.Wrap(errors.go:81)
    at github.com/BlockchainMiddleware/errors.TestExample(errors_test.go:86)
    at testing.tRunner(testing.go:865)
    at runtime.goexit(asm_amd64.s:1337)
cause by test
```



## Unwrap error object

errors support unwrap error object to get original cause error

```golang
err := fmt.Errorf("test")
err1 := errors.Wrap(err,"layer 1")
err2 := errors.Wrap(err,"layer 2")

// require.Equal(t,errors.Unwrap(err2),err)

```

## Enhance error interface

errors support create error object with vendor and error code attributes.

```golang
ec := errors.New("test",errors.WithCode(-2),errors.WithVendor("test"))

vendor,ok := errors.Vendor(ec)
// require.True(t,ok)
// require.Equal(t,vendor,"test")

code,ok := errors.Code(ec)
// require.True(t,ok)
// require.Equal(t,code,-2)

```

> one more thing: errors also support bind customer attributes. for details lookup [test example](./errors_test.go)
