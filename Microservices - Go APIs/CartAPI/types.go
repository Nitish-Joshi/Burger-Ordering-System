package main

import "gopkg.in/mgo.v2/bson"

type Product struct {
	ProductName string `bson:"ProductName" json:"ProductName"`
	ProductDesc string `bson:"ProductDesc" json:"ProductDesc"`
	Price       int    `bson:"Price" json:"Price"`
	Image       string `bson:"Image" json:"Image"`
	Quantity    int    `bson:"Quantity" json:"Quantity"`
}

type Cart struct {
	CartID   bson.ObjectId `bson:"_id" json:"CartID"`
	Total    int           `bson:"Total" json:"Total"`
	Products []Product     `bson:"Products" json:"Products"`
}
