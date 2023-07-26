package goerr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackTrace(t *testing.T) {
	u := &user{
		username: "Juanito valderrama",
	}
	s := userService{}
	err := s.createUser(u)
	assert.Error(t, err, "userService.createUser: userService.addRole: syntax error at near 'INSERT'")
}

func TestWrapGrpc(t *testing.T) {
	t.Run("should hide internal error message", func(t *testing.T) {
		u := &user{
			username: "Juanito valderrama",
		}
		s := userService{}
		err := s.createUser(u)
		err = WrapGrpc(err)
		assert.Error(t, err, "rpc error: code = Internal desc = Internal server error")
	})

	t.Run("should include error message", func(t *testing.T) {
		u := &user{
			username: "Juanito valderrama",
			age:      16,
		}
		s := userService{}
		err := s.validate(u)
		err = WrapGrpc(err)
		assert.Error(t, err, "rpc error: code = InvalidArgument desc = Age must be >= 18")
	})
}

type user struct {
	username string
	age      int
}

type userService struct{}

func (s *userService) createUser(user *user) error {
	if user.username == "" {
		return &Error{Code: Invalid, Message: "Username is required"}
	}
	if err := s.addRole(user, "default"); err != nil {
		return &Error{Op: "userService.createUser", Err: err}
	}

	return nil
}

func (s *userService) validate(user *user) error {
	if user.username == "" {
		return &Error{Code: Invalid, Message: "Username is required"}
	}
	if user.age < 18 {
		return &Error{Code: Invalid, Message: "Age must be >= 18"}
	}
	return nil
}

func (us *userService) addRole(user *user, role string) error {
	// simulate sql call
	sqlError := errors.New("syntax error at near 'INSERT'")
	return &Error{Op: "userService.addRole", Err: sqlError}
}
