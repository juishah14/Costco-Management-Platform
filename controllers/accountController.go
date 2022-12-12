package controller

// NEED TO CHANGE THIS FOR ACCOUNT

import (
	"Golang-Management-Platform/database"
	"Golang-Management-Platform/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var accountCollection *mongo.Collection = database.OpenCollection(database.Client, "account")

// Will implement this func later
func AccountsByMembership(id string) (Accounts []primitive.M, err error) {
	return Accounts, err
}

func GetAccountsByMembership() gin.HandlerFunc {
	return func(c *gin.Context) {
		membershipId := c.Param("membership_id")
		allAccounts, err := AccountsByMembership(membershipId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing accounts by their membership ID."})
			return
		}
		c.JSON(http.StatusOK, allAccounts)
	}
}

func GetAccounts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := accountCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing accounts."})
			return
		}

		var allAccounts []bson.M
		err = result.All(ctx, &allAccounts)
		if err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allAccounts)
	}
}

func GetAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var account models.Account
		accountId := c.Param("account_id")

		err := accountCollection.FindOne(ctx, bson.M{"account_id": accountId}).Decode(&account)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing this account."})
			return
		}
		c.JSON(http.StatusOK, account)
	}
}

func CreateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var account models.Account

		// BindJSON reads the body buffer (JSON from Postman) to deserialize it into a Golang struct
		err := c.BindJSON(&account)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Handle validation (validate the data based on our account struct)
		accountErr := validate.Struct(account)
		if accountErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": accountErr.Error()})
			return
		}

		// Finish creating the account object
		account.ID = primitive.NewObjectID()
		account.Account_id = account.ID.Hex()
		account.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		account.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// Insert this account into our account collection in Mongo DB
		result, insertErr := accountCollection.InsertOne(ctx, account)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, could not create this account."})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var account models.Account
		accountId := c.Param("account_id")

		// BindJSON reads the body buffer (JSON from Postman) to deserialize it into a Golang struct
		err := c.BindJSON(&account)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create an update object and append all necessary details to it
		var updateObj primitive.D

		if account.User_id != "" {
			updateObj = append(updateObj, bson.E{"user_id", account.User_id})
		}

		// Set update object's 'Updated_at' field
		account.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", account.Updated_at})

		// Update w MongoDB
		upsert := true
		filter := bson.M{"account_id": accountId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := accountCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, failed to update this account."})
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
