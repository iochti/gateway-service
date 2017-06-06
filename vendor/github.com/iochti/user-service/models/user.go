package models

import (
	"encoding/json"
	"time"
)

// User represents a platform user
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Login     string    `json:"login"`
	AvatarURL string    `json:"avatar_url"`
	Name      string    `json:"name"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

// ToByte converts the user in a byte array
func (u *User) ToByte() ([]byte, error) {
	ub, err := json.Marshal(u)
	if err != nil {
		return []byte{}, err
	}
	return ub, nil
}
