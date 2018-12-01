package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

// MongoDB Config
var mongodb_server = "mongodb"
var mongodb_database = "cmpe281"
var mongodb_collection = "burgerOrder"



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
	mx.HandleFunc("/list", userOdersHandler(formatter)).Methods("GET")
	mx.HandleFunc("/insert", burgerNewOrderHandler(formatter)).Methods("POST")
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")


}

// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API version 1.0 alive!"})
	}
}


//*********************************** Burger New Order Handler*********************
func burgerNewOrderHandler(formatter *render.Render) http.HandlerFunc{
	return func(w http.ResponseWriter, req *http.Request){
		session, err := mgo.Dial(mongodb_server)
		if err != nil{
			panic(err)
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)
		var result bson.M
		// result = bson.M{"OrderId" : "12345", "UserId" : "01111", "totalPrice" : "5.1"}
		req.ParseForm()
		x := req.Form.Get("OrderId")
		y := req.Form.Get("UserId")
		z := req.Form.Get("totalPrice")
		result = bson.M{"OrderId" : x, "UserId" : y, "totalPrice" : z}
		err = c.Insert(result)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println("Burger orders:", result )
		formatter.JSON(w, http.StatusOK, result)
	}
}




//*********************************** search user orders *********************

func userOdersHandler(formatter *render.Render) http.HandlerFunc{
	return func(w http.ResponseWriter, req *http.Request){
		session, err := mgo.Dial(mongodb_server)
		if err != nil{
			panic(err)
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(mongodb_database).C(mongodb_collection)
		var result[] bson.M
		req.ParseForm()                     // Parses the request body
		var uid string= req.Form.Get("UserId")
		err = c.Find(bson.M{"UserId" : uid}).All(&result)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println("Burger orders:", result )
		formatter.JSON(w, http.StatusOK, result)
	}
}












/*

  	-- Gumball MongoDB Collection (Create Document) --

    db.gumball.insert(
	    { 
	      Id: 1,
	      CountGumballs: NumberInt(202),
	      ModelNumber: 'M102988',
	      SerialNumber: '1234998871109' 
	    }
	) ;

    -- Gumball MongoDB Collection - Find Gumball Document --

    db.gumball.find( { Id: 1 } ) ;

    {
        "_id" : ObjectId("54741c01fa0bd1f1cdf71312"),
        "Id" : 1,
        "CountGumballs" : 202,
        "ModelNumber" : "M102988",
        "SerialNumber" : "1234998871109"
    }

    -- Gumball MongoDB Collection - Update Gumball Document --

    db.gumball.update( 
        { Dd: 1 }, 
        { $set : { CountGumballs : NumberInt(10) } },
        { multi : false } 
    )

    -- Gumball Delete Documents

    db.gumball.remove({})

 */
