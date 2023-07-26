package goerr

import (
	"bytes"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error defines standard domain/application error.
type Error struct {
	// Machine-readable error code.
	Code string
	// Human-readable message.
	Message string
	// Logical operation and nested error.
	Op  string
	Err error
}

// Machine-readble error codes.
const (
	Internal       = "internal"       // internal error
	Invalid        = "invalid"        // validation failed
	NotFound       = "not found"      // resource does not exist
	Authorization  = "authorization"  // insufficient or missing permissions
	Authentication = "authentication" // authentication required
)

// Error returns the string representation of the error message.
func (e *Error) Error() string {
	var buf bytes.Buffer
	// Print the current operation in our stack, if any.
	if e.Op != "" {
		fmt.Fprintf(&buf, "%s: ", e.Op)
	}
	// If wrapping an error, print its Error() message.
	// Otherwise print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			fmt.Fprintf(&buf, "<%s> ", e.Code)
		}
		buf.WriteString(e.Message)
	}
	return buf.String()
}

// ErrorCode returns the code of the root error, if exits. Otherwise returns Internal.
func ErrorCode(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Code != "" {
		return e.Code
	} else if ok && e.Err != nil {
		return ErrorCode(e.Err)
	}
	return Internal
}

// ErrorMessage returns the human-readable message of the error, if exists.
// Otherwise returns a generic error message.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Message != "" {
		return e.Message
	} else if ok && e.Err != nil {
		return ErrorMessage(e.Err)
	}
	return "An internal error has occurred"
}

// Transform domain Error to grpc error.
func WrapGrpc(err error) error {
	// Check if is Error
	e, ok := err.(*Error)
	if !ok {
		return status.Error(codes.Unknown, err.Error())
	}
	code := codes.Internal
	switch e.Code {
	case Invalid:
		code = codes.InvalidArgument
	case NotFound:
		code = codes.NotFound
	case Authorization:
		code = codes.PermissionDenied
	case Authentication:
		code = codes.Unauthenticated
	}
	msg := e.Error()
	// Internal error or no error code msg must be hide.
	if e.Code == Internal || e.Code == "" {
		msg = "Internal server error"
	}
	return status.Error(code, msg)
}
