// Recipes API
//
//	This is a sample recipes API. You can find out more about the API at https://github.com/richinex/recipes-api.
//
//	Schemes: http
//	Host: localhost:8080
//		BasePath: /
//		Version: 1.0.0
//		Contact: Richard Chukwu <richinex@gmail.com>
//	SecurityDefinitions:
//	api_key:
//	type: apiKey
//	name: Authorization
//	in: header
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

	"github.com/go-redis/redis"
	"github.com/richinex/recipes-api/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type application struct {
	recipesModel *models.RecipesModel
	usersModel   *models.UsersModel
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
	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	app := &application{
		recipesModel: &models.RecipesModel{
			Collection:  collection,
			Ctx:         ctx,
			RedisClient: redisClient,
		},
		usersModel: &models.UsersModel{
			Collection: collectionUsers,
		},
	}

	status := redisClient.Ping()
	log.Println(status)

	app.routes()

}
