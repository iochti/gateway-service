package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Thing represents a Db Thing element
type Thing struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	GroupID   string        `json:"group_id" bson:"groupID"`
	Name      string        `json:"name" bson:"name"`
	MAC       string        `json:"mac" bson:"mac"`
	CreatedAt time.Time     `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updatedAt"`
}
