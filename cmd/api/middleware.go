package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/richinex/recipes-api/internal/models"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
)

func (app *application) signInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the user from the database
	cur := app.usersModel.Collection.FindOne(app.usersModel.Ctx, bson.M{"username": user.Username})
	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Get the salt value and hashed password from the database
	var dbUser models.User
	err := cur.Decode(&dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}
	hashedPassword := dbUser.Password

	// Hash the password
	h := sha256.New()
	io.WriteString(h, user.Password)
	computedHash := hex.EncodeToString(h.Sum(nil))

	// Compare the computed hash with the hashed password from the database
	if computedHash != hashedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "User signed in"})
}

// swagger:operation POST /signout auth signOut
// Signing out
// ---
// responses:
//     '200':
//         description: Successful operation

func (app *application) signOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Signed out..."})
}

func (app *application) signUpHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user already exists
	count, err := app.usersModel.Collection.CountDocuments(app.usersModel.Ctx, bson.M{"username": user.Username})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash the user's password
	h := sha256.New()
	h.Write([]byte(user.Password))
	hashedPassword := hex.EncodeToString(h.Sum(nil))

	// Create the user object to insert
	newUser := bson.M{
		"username": user.Username,
		"password": hashedPassword,
	}

	// Insert the user into the database
	result, err := app.usersModel.Collection.InsertOne(app.usersModel.Ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}

	// Return the user ID
	c.JSON(http.StatusOK, gin.H{"user_id": result.InsertedID})
}

func (app *application) refreshHandler(c *gin.Context) {
	session := sessions.Default(c)
	sessionToken := session.Get("token")
	sessionUser := session.Get("username")
	if sessionToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session cookie"})
		return
	}

	sessionToken = xid.New().String()
	session.Set("username", sessionUser.(string))
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "New session issued"})
}

func (app *application) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Not logged",
			})
			c.Abort()
		}
		c.Next()
	}
}
