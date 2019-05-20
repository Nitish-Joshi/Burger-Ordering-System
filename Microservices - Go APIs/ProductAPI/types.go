package main

import "gopkg.in/mgo.v2/bson"

type Product struct {
	ProductID   bson.ObjectId `bson:"_id" json:"productID"`
	ProductName string        `bson:"ProductName" json:"ProductName"`
	ProductDesc string        `bson:"ProductDesc" json:"ProductDesc"`
	Price       string        `bson:"Price" json:"Price"`
	Image       string        `bson:"Image" json:"Image"`
}
