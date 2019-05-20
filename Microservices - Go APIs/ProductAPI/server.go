package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB Config
var mongodb_server = "mongodb://admin:cmpe281@10.0.1.39,10.0.1.190,10.0.3.176,11.0.1.200,11.0.1.247"

var mongodb_database = "TeamProject"
var mongodb_collection = "products"

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/products", productsHandler(formatter)).Methods("GET")
	mx.HandleFunc("/products/{productID}", getProductsHandler(formatter)).Methods("GET")
	mx.HandleFunc("/products", addProductHandler(formatter)).Methods("POST")
	mx.HandleFunc("/products/{productID}", deleteProductsHandler(formatter)).Methods("DELETE")
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API version 1.0 alive!"})
	}
}

// API  Handler --------------- Get all the products (GET) ------------------
func productsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var products []Product

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		if err = c.Find(bson.M{}).All(&products); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, products)
	}
}

// API  Handler --------------- Get the product info (GET) ------------------
func getProductsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		productId := vars["productID"]

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var result Product
		if err = c.FindId(bson.ObjectIdHex(productId)).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, result)
	}
}

// API  Handler --------------- Add a product (POST) ------------------
func addProductHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		fmt.Println("Connected to the database")

		var newProduct Product
		if err := json.NewDecoder(req.Body).Decode(&newProduct); err != nil {
			formatter.JSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		fmt.Println(newProduct)

		newProduct.ProductID = bson.NewObjectId()

		if err := c.Insert(&newProduct); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var result Product
		if err = c.Find(bson.M{"ProductName": newProduct.ProductName}).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, result)
	}
}

//API Handler --------------- Delete a product (DELETE) ------------------
func deleteProductsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		productID := vars["productID"]

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		if err = c.Remove(bson.M{"productID": productID}); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, "Product has been deleted successfully!!")
	}
}
