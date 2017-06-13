package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// ThingGroup represents a group of thing
type ThingGroup struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	AccountID   bson.ObjectId `json:"account_id" bson:"accountId"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	CreatedAt   time.Time     `json:"created_at" bson:"createdAt"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updatedAt"`
}
