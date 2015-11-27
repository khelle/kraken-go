package errors

import (
	"os"
	"fmt"
)

/**
 * Error interface
 */
type Error interface {
	GetCode()		int
	GetMessage()	string
	Error()			string
}

/**
 * Raise(Error)
 */
func Log(err Error) {
	if err != nil {
		fmt.Printf("Error[%d] = %s\n", err.GetCode(), err.GetMessage())
		os.Exit(err.GetCode())
	}
}

/**
 * ErrorInstance struct
 */
type ErrorInstance struct {
	code 	int
	message	string
}

/**
 * ErrorInstance constructor
 */
func New(code int, message string) *ErrorInstance {
	err := &ErrorInstance{}

	err.code = code
	err.message = message

	return err
}

/**
 * ErrorInstance.GetCode() int
 */
func (err *ErrorInstance) GetCode() int {
	return err.code
}

/**
 * ErrorInstance.GetMessage() string
 */
func (err *ErrorInstance) GetMessage() string {
	return err.message
}

/**
 * ErrorInstance.Error() string
 */
func (err *ErrorInstance) Error() string {
	return err.message
}
