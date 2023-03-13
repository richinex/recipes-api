package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *application) routes() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/recipes", app.listRecipesHandler)
	router.POST("/signin", app.signInHandler)
	router.POST("/signup", app.signUpHandler)
	router.POST("/refresh", app.refreshHandler)
	// router.POST("/signout", app.signOutHandler)

	authorized := router.Group("/")

	// authorized.Use(app.authMiddleware())
	// {
	authorized.POST("/recipes", app.newRecipeHandler)
	authorized.PUT("/recipes/:id", app.updateRecipeHandler)
	authorized.DELETE("/recipes/:id", app.deleteRecipeHandler)
	authorized.GET("/recipes/:id", app.getOneRecipeHandler)
	// }
	router.Run()
}
