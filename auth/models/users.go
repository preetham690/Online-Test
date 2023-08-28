package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	User_ID  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string
	Password string
	isAdmin  bool
}
