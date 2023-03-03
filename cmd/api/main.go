// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/richinex/recipes-api.
//
//		Schemes: http
//	 Host: localhost:8080
//		BasePath: /
//		Version: 1.0.0
//		Contact: Richard Chukwu <richinex@gmail.com>
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
// swagger:meta
package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type handlerApplication struct {
	recipesHandler *recipesHandler
}

var ctx context.Context
var err error
var client *mongo.Client

func main() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	appHandler := &handlerApplication{
		recipesHandler: newRecipesHandler(ctx, collection),
	}

	appHandler.routes()

}
