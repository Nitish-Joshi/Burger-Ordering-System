package main

import "gopkg.in/mgo.v2/bson"

type Payment struct {
	PaymentID     bson.ObjectId `bson:"_id" json:"PaymentID"`
	Amount        string        `bson:"Amount" json:"Amount"`
	PaymentStatus string        `bson:"PaymentStatus" json:"PaymentStatus"`
	UserID        string        `bson:"UserID" json:"UserID"`
	OrderID       string        `bson:"OrderID" json:"OrderID"`
}
