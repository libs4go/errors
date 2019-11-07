package errors

// ErrTODO errors
var ErrTODO = New("todo")

// TODO invoke todo panic
func TODO(message string) {
	panic(Wrap(ErrTODO, message))
}
