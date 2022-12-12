package controller

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

var membershipCollection *mongo.Collection = database.OpenCollection(database.Client, "membership")

func GetMemberships() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching membership details."})
		}

		var allMemberships []bson.M
		err = result.All(ctx, &allMemberships)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMemberships)
	}
}

func GetMembership() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var membership models.Membership
		membershipId := c.Param("membership_id")

		err := membershipCollection.FindOne(ctx, bson.M{"membership_id": membershipId}).Decode(&membership)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while fetching this membership's details."})
		}
		c.JSON(http.StatusOK, membership)
	}
}

func CreateMembership() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var membership models.Membership

		// BindJSON reads the body buffer (JSON from Postman) to deserialize it into a Golang struct
		err := c.BindJSON(&membership)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Handle validation (validate the data based on our membership struct)
		validationErr := validate.Struct(membership)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Finish creating the membership object
		membership.ID = primitive.NewObjectID()
		membership.Membership_id = membership.ID.Hex()
		membership.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		membership.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// Insert this membership into our membership collection in Mongo DB
		result, insertErr := membershipCollection.InsertOne(ctx, membership)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, could not create this membership."})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateMembership() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var membership models.Membership
		membershipId := c.Param("membership_id")

		// BindJSON reads the body buffer (JSON from Postman) to deserialize it into a Golang struct
		err := c.BindJSON(&membership)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create an update object and append all necessary details to it
		var updateObj primitive.D
		if membership.Membership_type != "" {
			updateObj = append(updateObj, bson.E{"membership_type", membership.Membership_type})
		}

		// Set update object's 'Updated_at' field
		membership.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// Update w MongoDB
		upsert := true
		filter := bson.M{"membership_id": membershipId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := membershipCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, failed to update this membership."})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func AccountMembershipCreator(membership models.Membership) string {

	// Create a new membership for a new account
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	membership.ID = primitive.NewObjectID()
	membership.Membership_id = membership.ID.Hex()
	membership.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	membership.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	// Insert the membership into our membership collection
	membershipCollection.InsertOne(ctx, membership)
	defer cancel()

	// Return the new membership's membership_id
	return membership.Membership_id
}
