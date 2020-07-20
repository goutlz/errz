package errz

import (
	"errors"
	"fmt"
	"runtime"
)

var (
	captureStacks = false
)

func SetStackCaptureMode(enabled bool) {
	captureStacks = enabled
}

type stackInfo struct {
	line int
	file string
}

type Error error

type errorImpl struct {
	msg    string
	subErr Error
	stack  *stackInfo
}

func (c *errorImpl) string() string {
	if c.stack == nil {
		return fmt.Sprintf("\tError: %v\n\n", c.msg)
	}

	return fmt.Sprintf("%v:%v\n\tError: %v\n\n", c.stack.file, c.stack.line, c.msg)
}

func (c *errorImpl) Error() string {
	err := c.string()
	if c.subErr == nil {
		return err
	}

	cErr := &errorImpl{}
	if errors.As(c.subErr, &cErr) {
		return err + c.subErr.Error()
	}

	return err + fmt.Sprintf("\tError: %v\n\n", c.subErr)
}

func (c *errorImpl) Unwrap() error {
	return c.subErr
}

func (c *errorImpl) Is(target error) bool {
	cErr := &errorImpl{}
	if errors.As(target, &cErr) {
		return c.msg == cErr.msg
	}

	return c.msg == target.Error()
}

func captureStack(skip int) *stackInfo {
	if !captureStacks {
		return nil
	}

	stack := &stackInfo{}
	_, stack.file, stack.line, _ = runtime.Caller(skip)
	return stack
}

func newError(subErr Error, msg string, skip int) Error {
	return &errorImpl{
		subErr: subErr,
		msg:    msg,
		stack:  captureStack(skip),
	}
}

func New(msg string) Error {
	return newError(nil, msg, 3)
}

func Newf(format string, args ...interface{}) Error {
	return newError(nil, fmt.Sprintf(format, args...), 3)
}

func Wrap(subErr error, msg string) Error {
	return newError(subErr, msg, 3)
}

func Wrapf(subErr error, format string, args ...interface{}) Error {
	return newError(subErr, fmt.Sprintf(format, args...), 3)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
