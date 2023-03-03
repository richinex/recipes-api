package main

import "github.com/gin-gonic/gin"

func (appHandler *handlerApplication) routes() {
	router := gin.Default()
	router.POST("/recipes", appHandler.newRecipeHandler)
	router.GET("/recipes", appHandler.listRecipesHandler)
	router.PUT("/recipes/:id", appHandler.updateRecipeHandler)
	router.DELETE("/recipes/:id", appHandler.deleteRecipeHandler)
	router.GET("/recipes/:id", appHandler.getOneRecipeHandler)
	router.Run()
}
