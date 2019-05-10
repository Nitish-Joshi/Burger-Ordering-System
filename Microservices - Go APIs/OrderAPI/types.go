package main

import "gopkg.in/mgo.v2/bson"

type Product struct {
	ProductName string `bson:"ProductName" json:"ProductName"`
	ProductDesc string `bson:"ProductDesc" json:"ProductDesc"`
	Price       string `bson:"Price" json:"Price"`
	Image       string `bson:"Image" json:"Image"`
	Quantity    string `bson:"Quantity" json:"Quantity"`
}

type Order struct {
	OrderID     bson.ObjectId `bson:"_id" json:"OrderID"`
	UserID      string        `bson:"UserID" json:"UserID"`
	Total       string        `bson:"Total" json:"Total"`
	Products    []Product     `bson:"Products" json:"Products"`
	OrderStatus string        `bson:"OrderStatus" json:"OrderStatus"`
	PaymentStatus string      `bson:"PaymentStatus" json:"PaymentStatus"`
}
