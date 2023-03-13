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
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/richinex/recipes-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

type claims struct {
	Username           string `json:"username"`
	jwt.StandardClaims `json:"security,omitempty"`
}
type jwtOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

// swagger:operation POST /signin auth signIn
// Login with username and password
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//		description: Successful operation
//	'401':
//
// description: Invalid credentials
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

	// Generate a JWT token
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jwtOutput := jwtOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtOutput)
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

// swagger:operation POST /refresh auth refresh
// Get new token in exchange for an old one
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//		description: Successful operation
//	'400':
//		description: Token is new and doesn't need
//					 a refresh
//	'401':
//		description: Invalid credentials
func (app *application) refreshHandler(c *gin.Context) {
	tokenValue := c.GetHeader("Authorization")
	claims := &claims{}
	tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Updated code to use time.until
	if time.Until(time.Unix(claims.ExpiresAt, 0)) > 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is not expired yet"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jwtOutput := jwtOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	c.JSON(http.StatusOK, jwtOutput)
}

// func (app *application) authMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenValue := c.GetHeader("Authorization")
// 		claims := &claims{}
// 		tkn, err := jwt.ParseWithClaims(tokenValue, claims,
// 			func(token *jwt.Token) (interface{}, error) {
// 				return []byte(os.Getenv("JWT_SECRET")), nil
// 			})
// 		if err != nil {
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 		}
// 		if tkn == nil || !tkn.Valid {
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 		}
// 		c.Next()
// 	}
// }
