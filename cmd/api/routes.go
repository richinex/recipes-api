package main

import (
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func (app *application) routes() {
	router := gin.Default()
	store, _ := redis.NewStore(10, "tcp", os.Getenv("REDIS_URI"), os.Getenv("REDIS_PASSWORD"), []byte(os.Getenv("REDIS_KEY")))
	router.Use(sessions.Sessions("recipes_api", store))
	router.GET("/recipes", app.listRecipesHandler)
	router.POST("/signin", app.signInHandler)
	router.POST("/signup", app.signUpHandler)
	router.POST("/refresh", app.refreshHandler)
	router.POST("/signout", app.signOutHandler)

	authorized := router.Group("/")
	authorized.Use(app.authMiddleware())
	{
		authorized.POST("/recipes", app.newRecipeHandler)
		authorized.PUT("/recipes/:id", app.updateRecipeHandler)
		authorized.DELETE("/recipes/:id", app.deleteRecipeHandler)
		authorized.GET("/recipes/:id", app.getOneRecipeHandler)
	}
	router.Run()
}
