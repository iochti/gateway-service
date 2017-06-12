package models

import (
	"encoding/json"
	"time"
)

// User represents a platform user
type User struct {
	ID        string    `json:"id" bson:"id"`
	Email     string    `json:"email" bson:"email"`
	Login     string    `json:"login" bson:"login"`
	AvatarURL string    `json:"avatar_url" bson:"avatarURL"`
	Name      string    `json:"name" bson:"name"`
	Created   time.Time `json:"created" bson:"created"`
	Updated   time.Time `json:"updated" bson:"updated"`
}

// ToByte converts the user in a byte array
func (u *User) ToByte() ([]byte, error) {
	ub, err := json.Marshal(u)
	if err != nil {
		return []byte{}, err
	}
	return ub, nil
}
