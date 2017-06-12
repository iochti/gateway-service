package models

import "time"

// Thing represents a Db Thing element
type Thing struct {
	ID        string    `json:"id" bson:"id"`
	GroupID   string    `json:"group_id" bson:"groupID"`
	Name      string    `json:"name" bson:"name"`
	MAC       string    `json:"mac" bson:"mac"`
	CreatedAt time.Time `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time `json:"updated_at" bson:"updatedAt"`
}
