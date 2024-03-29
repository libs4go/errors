package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
)

// Debug print stack information flag
var Debug = true

type callStackError struct {
	callStack *callStack
	err       error
	message   string
}

func newCallStackError(message string, err error, skip int) *callStackError {

	callStack := newCallStack(skip)

	return &callStackError{
		err:       err,
		callStack: callStack,
		message:   message,
	}
}

func (err *callStackError) Error() string {

	var buff bytes.Buffer

	buff.WriteString(fmt.Sprintf("error: %s\n", err.message))

	if Debug {
		buff.WriteString(fmt.Sprintf("%s", err.callStack))
	}

	if err.err != nil {
		buff.WriteString(fmt.Sprintf("cause by %s", err.err))
	}

	return buff.String()
}

// StackTrace error raise stack trace function
func StackTrace(err error, tracer func(runtime.Frame)) {
	if csError, ok := err.(*callStackError); ok {
		csError.callStack.WalkCallFrames(tracer)
		return
	}
	// no callstack
	newCallStack(1).WalkCallFrames(tracer)
}

// Cause get error's cause error
func Cause(err error) error {
	if csError, ok := err.(*callStackError); ok {
		return csError.err
	}

	return nil
}

// Unwrap walk throught the error cause list and return the header cause error
func Unwrap(err error) error {
	for {
		cause := Cause(err)

		if cause == nil {
			return err
		}

		err = cause
	}
}

// Is Is unwraps its first argument sequentially looking for an error that matches the second.
// It reports whether it finds a match. It should be used in preference to simple equality checks:
func Is(first error, second error) bool {
	return Unwrap(first) == second
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool {

	if err == nil {
		return false
	}

	targetType := reflect.TypeOf(target)

	targetValue := reflect.ValueOf(target)

	if targetType.Kind() != reflect.Ptr || targetValue.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}

	targetType = targetType.Elem()

	if targetType.Kind() != reflect.Ptr {
		return false
	}

	if !targetType.Implements(errorInterfaceType) {
		panic("target must implement error interface")
	}

	unwraped := Unwrap(err)

	errType := reflect.TypeOf(unwraped)

	if errType != targetType {
		return false
	}

	targetValue.Elem().Set(reflect.ValueOf(unwraped))

	return true
}

var errorInterfaceType = reflect.TypeOf((*error)(nil)).Elem()

// Wrap wrap error with stacktrace
func Wrap(err error, fmtstr string, args ...interface{}) error {
	return newCallStackError(fmt.Sprintf(fmtstr, args...), err, 2)
}

// ErrorCode .
type ErrorCode struct {
	Vendor  string                 `json:"vendor"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Attrs   map[string]interface{} `json:"attrs"`
}

func (ec *ErrorCode) Error() string {

	return fmt.Sprintf("(%s:%d) %s", ec.Vendor, ec.Code, ec.Message)
}

// func (ec *ErrorCode) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(ec)
// }

// func (ec *ErrorCode) UnmarshalJSON(b []byte) error {
// 	return json.Unmarshal(b, ec)
// }

// Option .
type Option func(ec *ErrorCode)

// WithCode error code option
func WithCode(code int) Option {
	return func(ec *ErrorCode) {
		ec.Code = code
	}
}

// WithVendor error vendor option
func WithVendor(vendor string) Option {
	return func(ec *ErrorCode) {
		ec.Vendor = vendor
	}
}

// WithAttr bind error customer attribute
func WithAttr(name string, value interface{}) Option {
	return func(ec *ErrorCode) {
		if ec.Attrs == nil {
			ec.Attrs = make(map[string]interface{})
		}

		ec.Attrs[name] = value
	}
}

// FromJSON unmarshal errorcode from json
// if any error occur， return nil
func FromJSON(data []byte) error {

	var ec ErrorCode

	err := json.Unmarshal(data, &ec)

	if err != nil {
		return nil
	}

	return &ec
}

// New create errors enhance error object which support errorcode and vendor id
func New(errmsg string, options ...Option) error {

	ec := &ErrorCode{
		Message: errmsg,
		Code:    -1,
		Vendor:  "errors",
	}

	for _, option := range options {
		option(ec)
	}

	return ec
}

// Vendor get error associate vendor name
func Vendor(err error) (string, bool) {
	if ec, ok := Unwrap(err).(*ErrorCode); ok {
		return ec.Vendor, true
	}

	return "", false
}

// Code get error associate code
func Code(err error) (int, bool) {
	if ec, ok := Unwrap(err).(*ErrorCode); ok {
		return ec.Code, true
	}

	return -1, false
}

// Attr get associate attribute value
func Attr(err error, name string, value interface{}) bool {

	if ec, ok := Unwrap(err).(*ErrorCode); ok {

		if ec.Attrs == nil {
			return false
		}

		if attr, ok := ec.Attrs[name]; ok {

			valType := reflect.TypeOf(value)
			attrType := reflect.TypeOf(attr)

			if valType.Kind() != reflect.Ptr {
				return false
			}

			valType = valType.Elem()

			if valType.Kind() == reflect.Interface &&
				attrType.Kind() == reflect.Ptr &&
				attrType.Elem().Kind() == reflect.Struct &&
				attrType.Implements(valType) {

				reflect.ValueOf(value).Elem().Set(reflect.ValueOf(attr))

				return true
			}

			if valType == attrType {
				reflect.ValueOf(value).Elem().Set(reflect.ValueOf(attr))
				return true
			}
		}
	}

	return false
}
