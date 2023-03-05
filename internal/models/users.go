package models

import (
	"context"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type UsersModel struct {
	Collection  *mongo.Collection
	Ctx         context.Context
	RedisClient *redis.Client
}
