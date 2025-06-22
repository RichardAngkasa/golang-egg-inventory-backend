package main

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type EggTypeEnum string

const (
	EggTypeBig    EggTypeEnum = "big"
	EggTypeMedium EggTypeEnum = "medium"
	EggTypeSmall  EggTypeEnum = "small"
)

var validEggTypes = map[EggTypeEnum]bool{
	EggTypeBig:    true,
	EggTypeMedium: true,
	EggTypeSmall:  true,
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type EggRack struct {
	ID          uuid.UUID   `json:"id"`
	EggType     EggTypeEnum `json:"egg_type"`
	DateCreated time.Time   `json:"date_created"`
	User        string      `json:"User"`
}

func (e *EggRack) Validate() error {
	if !validEggTypes[e.EggType] {
		return errors.New("invalid egg type: must big, medium, or small")
	}
	if e.User == "" {
		return errors.New("user can't be empty")
	}
	return nil
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.New("username cannot be empty")
	}
	if u.Password == "" {
		return errors.New("password cannot be empty")
	}
	return nil
}
