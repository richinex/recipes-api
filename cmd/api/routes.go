package main

import "github.com/gin-gonic/gin"

func (app *application) routes() {
	router := gin.Default()
	router.POST("/recipes", app.newRecipeHandler)
	router.GET("/recipes", app.listRecipesHandler)
	router.PUT("/recipes/:id", app.updateRecipeHandler)
	router.DELETE("/recipes/:id", app.deleteRecipeHandler)
	router.GET("/recipes/:id", app.getOneRecipeHandler)
	router.Run()
}
