package main

import "gopkg.in/mgo.v2/bson"

type user struct {
	UserID   bson.ObjectId `bson:"_id" json:"UserID"`
	Name     string        `bson:"Name" json:"Name"`
	Email    string        `bson:"Email" json:"Email"`
	Address  string        `bson:"Address" json:"Address"`
	Password string        `bson:"Password" json:"Password"`
	CartID   string        `bson:"CartID" json:"CartID"`
}
