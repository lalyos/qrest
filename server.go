// Command qrest is a quick RESTful JSON server
//
// How to use
//
// Create a JSON file containing the data you'd like to be part of your server. An example file might look like:
//
//    {
//        "posts": [ { "id": 1, "title": "Foo" } ]
//    }
//
// Start qrest with this file as an argument:
//
//    qrest db.json
//
// Or in a docker container:
//
//    $ docker build -t qrest .
//    $ docker run --rm -p 3000:3000 qrest "db.json" # assuming db.json is in this source directory
//
// This will create the following routes for you to use:
//
//    POST /posts (creates a new post record)
//    GET /posts (returns all post records)
//    GET /posts/:id (returns a specific record)
//    PUT /posts/:id (creates or updates a record with the specified ID)
//    PATCH /posts/:id (updates a record with the specified ID)
//    DELETE /posts/:id (deletes the specified record)
//
//
package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	nlogrus "github.com/meatballhat/negroni-logrus"
)

var (
	logger *logrus.Logger
)

func main() {
	logr := nlogrus.NewMiddleware()
	logger = logr.Logger

	parseJsonFile()

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":3000"
	}

	router := httprouter.New()

	addStaticRoutes(router)
	addDynamicRoutes(router)

	// This goroutine will flush the JSON to the db.json file every 30 seconds,
	// OR before the application exits
	go flushJson()

	n := negroni.Classic()
	n.Use(logr)
	n.UseHandler(router)
	n.Run(port)
}
