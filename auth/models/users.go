package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//	type User struct {
//		User_ID  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
//		Username string
//		Password string
//		isAdmin  bool
//	}
type User struct {
	User_ID  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"` //this is the login ID
	Mobile   string             `json:"mobile" bson:"mobile"`
	Password string             `json:"password" bson:"password"` //and this is the password
	IsAdmin  bool               `json:"isAdmin" bson:"isAdmin"`
	//LastModified int64              `json:"lastModified" bson:"lastModified"`
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("invalid email address field")
	}
	if u.Password == "" {
		//return errors.New("invalid email field")
		return fmt.Errorf("invalid password field")
	}
	if u.Mobile == "" {
		return fmt.Errorf("invalid mobile number")
	}
	return nil
}

func (u *User) ToBytes() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) ToString() (string, error) {
	bytes, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
