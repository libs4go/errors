package errors

import "testing"

func TestCallStack(t *testing.T) {
	func() {
		func() {
			func() {
				println(newCallStack(1).String())
			}()
		}()
	}()
}
