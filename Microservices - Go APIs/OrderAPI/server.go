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
var mongodb_server = "mongodb://admin:cmpe281@10.0.1.249,10.0.1.9,10.0.1.84,11.0.1.237,11.0.1.54"

//var mongodb_server1 string
//var mongodb_server2 string
//var redis_server string

//="mongodb://54.67.13.87:27017,54.67.106.101:27017,13.57.39.192:27017,54.153.26.217://27017,52.53.154.42:27017"

var mongodb_database = "TeamProject"
var mongodb_collection = "order"

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
	mx.HandleFunc("/orderofusers/{userid}", getAllOrdersHandler(formatter)).Methods("GET")
	mx.HandleFunc("/orders/{orderid}", orderHandler(formatter)).Methods("GET")
	mx.HandleFunc("/orders", newOrderHandler(formatter)).Methods("POST")
	mx.HandleFunc("/orders/{orderid}", deleteOrderHandler(formatter)).Methods("DELETE")
	mx.HandleFunc("/orders", updateOrderHandler(formatter)).Methods("PUT")
	mx.HandleFunc("/orders/updateorderstatus/{orderid}", updateOrderStatusHandler(formatter)).Methods("PUT")
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API version 1.0 alive!"})
	}
}

// API  Handler --------------- Get all Orders (GET) ------------------
func getAllOrdersHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		userid := vars["userid"]

		var orders []Order

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		if err = c.Find(bson.M{"UserID": userid}).All(&orders); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, orders)
	}
}

// API  Handler --------------- Get the order info (GET) ------------------
func orderHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		orderid := vars["orderid"]

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var result Order
		if err = c.FindId(bson.ObjectIdHex(orderid)).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, result)
	}
}

// API  Handler --------------- Register a new order (POST) ------------------
func newOrderHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var newOrder Order
		if err := json.NewDecoder(req.Body).Decode(&newOrder); err != nil {
			formatter.JSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		fmt.Println(newOrder)

		newOrder.OrderID = bson.NewObjectId()
		newOrder.OrderStatus = "Order has been placed!!"
		newOrder.PaymentStatus = "Pending"

		if err := c.Insert(&newOrder); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var result Order
		// if err = c.FindId(bson.ObjectIdHex(order.OrderID)).One(&result); err != nil {
		// 	formatter.JSON(w, http.StatusInternalServerError, err.Error())
		// 	return
		// }

		if err = c.FindId(newOrder.OrderID).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, result)
	}
}

// API  Handler --------------- Delete the order (DELETE) ------------------
func deleteOrderHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		orderid := vars["orderid"]

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		fmt.Println(orderid)

		if err := c.Remove(bson.M{"OrderID": orderid}); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, "Order has been deleted successfully!!")
	}
}

// API  Handler --------------- Update the order (PUT) ------------------
func updateOrderHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var newOrder Order
		if err := json.NewDecoder(req.Body).Decode(&newOrder); err != nil {
			formatter.JSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if err := c.UpdateId(newOrder.OrderID, &newOrder); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var updatedOrder Order

		if err = c.FindId(newOrder.OrderID).One(&updatedOrder); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		//c.Find(bson.M{"newOrder.OrderID": newOrder.OrderID}).One(&updatedOrder)
		formatter.JSON(w, http.StatusOK, updatedOrder)

	}
}

// API  Handler --------------- Update the order status (PUT) ------------------
func updateOrderStatusHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		orderid := vars["orderid"]

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var result Order
		if err = c.FindId(bson.ObjectIdHex(orderid)).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		result.OrderStatus = "Order has been dispatched!!"
		result.PaymentStatus = "Payment received"

		if err := c.UpdateId(result.OrderID, &result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var updatedOrder Order

		if err = c.FindId(result.OrderID).One(&updatedOrder); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, updatedOrder)

	}
}

