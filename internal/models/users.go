package models

import (
	"context"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

// API user credentials
// It is used to sign in
//
// swagger:model user
type User struct {
	// User's password
	//
	// required: true
	Password string `json:"password"`
	// User's login
	//
	// required: true
	Username string `json:"username"`
}

// swagger:model userModel
type UsersModel struct {
	Collection  *mongo.Collection
	Ctx         context.Context
	RedisClient *redis.Client
}
