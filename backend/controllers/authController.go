package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/ostd2000/Auth/database"
	"github.com/ostd2000/Auth/jwt"
	"github.com/ostd2000/Auth/models"
)

var client *mongo.Client = database.DBInstance()
var userCollection *mongo.Collection = database.OpenCollection(client, "user") 
var validate = validator.New()

func HashPassword(password string) string {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(hashedPassword)
}

// "userPassword" is the hashed password stored in the database.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))

	check := true
	msg := ""
	if err != nil {
	  msg = fmt.Sprint("Email or password is incorrect.")	
		check = false
	}

	return check, msg
}

func AuthSuccessResponse(data *mongo.InsertOneResult) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data": data, 
		"error": nil,
	}
}

func AuthErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"data": "",
		"error": err.Error(),
	}
}

func Signup() fiber.Handler {
    return func(c *fiber.Ctx) error {
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
        var user models.User

		err := c.BodyParser(&user)
		if err != nil {
			c.Status(http.StatusBadRequest)

            return c.JSON(AuthErrorResponse(err))
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.Status(http.StatusBadRequest)

			return c.JSON(AuthErrorResponse(errors.New(
				"Validation error.",
			)))
		}

		filter := bson.M{"username": user.Username}
		count, err := userCollection.CountDocuments(ctx, filter)
		defer cancel()
		
		if err != nil {
			log.Panic(err)

			c.Status(http.StatusInternalServerError)

			return c.JSON(AuthErrorResponse(errors.New(
				"Error checking for the username.",
			)))
		}

		if count != 0 {
			c.Status(http.StatusInternalServerError)

			return c.JSON(AuthErrorResponse(errors.New(
				"username already exists.",
			))) 
		}

		hashedPassword := HashPassword(*user.Password)
		user.Password = &hashedPassword

        accessToken, err := jwt.GenerateAccessToken(user.UserID) 
	    refreshToken, err := jwt.GenerateRefreshToken(user.UserID)

	    user.AccessToken = &accessToken
	    user.RefreshToken = &refreshToken
	  
	    user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	    user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	    result, insertErr := userCollection.InsertOne(ctx, user)
		defer cancel()

	    if insertErr != nil {
			msg := fmt.Sprint("User item was not created.")

			c.Status(http.StatusInternalServerError)

		    return c.JSON(AuthErrorResponse(errors.New(msg)))
		}

		c.Status(http.StatusOK)

	    return c.JSON(AuthSuccessResponse(result))
	}
}

func Login() {}

func Logout() {}

func ForgetPassword() {}




















