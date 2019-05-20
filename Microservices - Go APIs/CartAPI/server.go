package main

import (
	"encoding/json"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB Config
var mongodb_server = "mongodb://admin:cmpe281@11.0.1.171,11.0.1.96,11.0.1.193,12.0.1.209,12.0.1.140"

//var mongodb_server1 string
//var mongodb_server2 string
//var redis_server string

//="mongodb://54.67.13.87:27017,54.67.106.101:27017,13.57.39.192:27017,54.153.26.217://27017,52.53.154.42:27017"

var mongodb_database = "TeamProject"
var mongodb_collection = "cart"

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	//mongodb_server = os.Getenv("MONGO1")
	//mongodb_server1 = os.Getenv("MONGO2")
	//mongodb_server2 = os.Getenv("MONGO3")
	//mongodb_database = os.Getenv("MONGO_DB")
	//mongodb_collection = os.Getenv("MONGO_COLLECTION")
	//redis_server = os.Getenv("REDIS")

	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
	//mx.HandleFunc("/users", getAllUsersHandler(formatter)).Methods("GET")
	mx.HandleFunc("/carts/{cartid}", cartHandler(formatter)).Methods("GET")
	mx.HandleFunc("/carts", createCartHandler(formatter)).Methods("POST")
	//mx.HandleFunc("/users/{email}", deleteUserHandler(formatter)).Methods("DELETE")
	mx.HandleFunc("/carts", updateCartHandler(formatter)).Methods("PUT")

}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API version 1.0 alive!"})
	}
}

// API  Handler --------------- Get the cart info (GET) ------------------
func cartHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		cartid := vars["cartid"]

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var result Cart
		if err = c.FindId(bson.ObjectIdHex(cartid)).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, result)
	}
}

// API  Handler --------------- Register a new cart (POST) ------------------
func createCartHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var newCart Cart
		if err := json.NewDecoder(req.Body).Decode(&newCart); err != nil {
			formatter.JSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		newCart.CartID = bson.NewObjectId()

		if err := c.Insert(&newCart); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var result Cart
		if err = c.FindId(newCart.CartID).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, result)
	}
}

// API  Handler --------------- Update the cart (PUT) ------------------
func updateCartHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var newCart Cart
		if err := json.NewDecoder(req.Body).Decode(&newCart); err != nil {
			formatter.JSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if err := c.UpdateId(newCart.CartID, &newCart); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var updatedCart Cart
		if err := c.FindId(newCart.CartID).One(&updatedCart); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, updatedCart)

	}
}
