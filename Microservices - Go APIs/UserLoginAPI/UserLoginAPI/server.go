package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB Config
var mongodb_server = os.Getenv("MONGO_SERVERS")
var mongodb_database = os.Getenv("MONGO_DB")
var mongodb_collection = os.Getenv("MONGO_COLLECTION")

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
	mx.HandleFunc("/users", getAllUsersHandler(formatter)).Methods("GET")
	mx.HandleFunc("/users/{email}", userHandler(formatter)).Methods("GET")
	mx.HandleFunc("/users", userSignUpHandler(formatter)).Methods("POST")
	mx.HandleFunc("/users/{email}", deleteUserHandler(formatter)).Methods("DELETE")
	mx.HandleFunc("/users", updateUserHandler(formatter)).Methods("PUT")

}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API version 1.0 alive!"})
	}
}

// API  Handler --------------- Get all Users (GET) ------------------
func getAllUsersHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var users []user

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		if err = c.Find(bson.M{}).All(&users); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, users)
	}
}

// API  Handler --------------- Get the user info (GET) ------------------
func userHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		email := vars["email"]

		session, err := mgo.Dial(mongodb_server)
		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var result user
		err = c.Find(bson.M{"Email": email}).One(&result)

		formatter.JSON(w, http.StatusOK, result)
	}
}

// API  Handler --------------- Register a new user (POST) ------------------
func userSignUpHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var newUser user
		if err := json.NewDecoder(req.Body).Decode(&newUser); err != nil {
			formatter.JSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		fmt.Println(newUser)

		newUser.UserID = bson.NewObjectId()

		if err := c.Insert(&newUser); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var result user
		if err = c.Find(bson.M{"Email": newUser.Email}).One(&result); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, result)
	}
}

// API  Handler --------------- Delete the user (DELETE) ------------------
func deleteUserHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		email := vars["email"]

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		if err := c.Remove(bson.M{"Email": email}); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		formatter.JSON(w, http.StatusOK, "User has been deleted successfully!!")
	}
}

// API  Handler --------------- Update the user (PUT) ------------------
func updateUserHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		session, err := mgo.Dial(mongodb_server)

		if err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer session.Close()
		session.SetMode(mgo.PrimaryPreferred, true)
		c := session.DB(mongodb_database).C(mongodb_collection)

		var newUser user
		if err := json.NewDecoder(req.Body).Decode(&newUser); err != nil {
			formatter.JSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if err := c.UpdateId(newUser.UserID, &newUser); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		var updatedUser user
		c.Find(bson.M{"Email": newUser.Email}).One(&updatedUser)
		formatter.JSON(w, http.StatusOK, updatedUser)

	}
}
