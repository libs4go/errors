package errors

// ErrTODO errors
var ErrTODO = New("todo")

// TODO invoke todo panic
func TODO(message string) error {
	panic(Wrap(ErrTODO, message))
}
