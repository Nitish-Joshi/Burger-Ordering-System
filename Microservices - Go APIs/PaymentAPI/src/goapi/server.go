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
var mongodb_server = "mongodb://admin:admin@10.0.2.249,10.0.1.16,10.0.2.116,10.1.2.197,10.1.1.207"
var mongodb_database = "cmpe281"
var mongodb_collection = "payment"

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
	mx.HandleFunc("/payments", paymentsHandler(formatter)).Methods("GET")
	mx.HandleFunc("/payments/{paymentid}", getPaymentHandler(formatter)).Methods("GET")
	mx.HandleFunc("/payments", addNewPaymentHandler(formatter)).Methods("POST")
	mx.HandleFunc("/payments/updateThePaymentStatus/{paymentid}", updatePaymentStatusHandler(formatter)).Methods("PUT")
	mx.HandleFunc("/payments/{paymentid}", deletePaymentHandler(formatter)).Methods("DELETE")
	mx.HandleFunc("/payments", updatePaymentHandler(formatter)).Methods("PUT")
	mx.HandleFunc("/payments/paymentfromorder/{orderid}", getPaymentFromOrderIDHandler(formatter)).Methods("GET")
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API version 1.0 alive!"})
	}
}

// API  Handler --------------- Get all the payments (GET) ------------------
func paymentsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var payments []Payment

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		if err = c.Find(bson.M{}).All(&payments); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, payments)
	}
}

// API  Handler --------------- Get the payment info (GET) ------------------
func getPaymentHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		paymentid := vars["paymentid"]

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var result Payment
		if err = c.FindId(bson.ObjectIdHex(paymentid)).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, result)
	}
}

// API  Handler --------------- Add a payment (POST) ------------------
func addNewPaymentHandler(formatter *render.Render) http.HandlerFunc {
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

		var newPayment Payment
		if err := json.NewDecoder(req.Body).Decode(&newPayment); err != nil {
			formatter.JSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		fmt.Println(newPayment)

		newPayment.PaymentID = bson.NewObjectId()

		if err := c.Insert(&newPayment); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var result Payment
		if err = c.FindId(newPayment.PaymentID).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, result)
	}
}

//API Handler --------------- Delete a payment (DELETE) ------------------
func deletePaymentHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		paymentid := vars["paymentid"]

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		if err := c.RemoveId(bson.ObjectIdHex(paymentid)); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, "Payment has been deleted successfully!!")
	}
}

// API  Handler --------------- Update the order (PUT) ------------------
func updatePaymentHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var newpayment Payment
		if err := json.NewDecoder(req.Body).Decode(&newpayment); err != nil {
			formatter.JSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if err := c.UpdateId(newpayment.PaymentID, &newpayment); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var updatedPayment Payment

		if err = c.FindId(newpayment.PaymentID).One(&updatedPayment); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, updatedPayment)

	}
}

// API  Handler --------------- Update the order (PUT) ------------------
func updatePaymentStatusHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		paymentid := vars["paymentid"]

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var result Payment
		if err = c.FindId(bson.ObjectIdHex(paymentid)).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		result.PaymentStatus = "Payment Received!!"

		if err := c.UpdateId(result.PaymentID, &result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var updatedPayment Payment

		if err = c.FindId(result.PaymentID).One(&updatedPayment); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, updatedPayment)

	}
}

// API  Handler --------------- Get the payment info (GET) ------------------
func getPaymentFromOrderIDHandler(formatter *render.Render) http.HandlerFunc {
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

		var result Payment
		if err = c.Find(bson.M{"OrderID": orderid}).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, result)
	}
}
