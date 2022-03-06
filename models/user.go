package models

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Hash     string `json:"-"`
	ID       int    `json:"-"`
}

type UserAuth struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

func (u *User) Validate() bool {
	u.Login = strings.TrimSpace(u.Login)
	u.Password = strings.TrimSpace(u.Password)
	return u.Login != "" && u.Password != ""
}

func ReadUserFromBody(r io.Reader) (*User, error) {
	user := &User{}

	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	if err := dec.Decode(user); err != nil {
		return nil, errors.New("Bad request: " + err.Error())
	}

	if !user.Validate() {
		return nil, errors.New("Bad request: login and/or password cannot be empty")
	}

	return user, nil
}
