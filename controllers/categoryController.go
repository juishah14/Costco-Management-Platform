package controller

// used to be menu

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

var categoryCollection *mongo.Collection = database.OpenCollection(database.Client, "category")

func GetCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := categoryCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing categories."})
		}

		var allCategories []bson.M
		if err = result.All(ctx, &allCategories); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allCategories)
	}
}

func GetCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var category models.Category
		categoryId := c.Param("category_id")

		err := categoryCollection.FindOne(ctx, bson.M{"category_id": categoryId}).Decode(&category)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing this category."})
		}
		c.JSON(http.StatusOK, category)
	}
}

func CreateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var category models.Category

		// BindJSON reads the body buffer (JSON from Postman) to deserialize it into a Golang struct
		err := c.BindJSON(&category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Handle validation (validate the data based on our category struct)
		validationErr := validate.Struct(category)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Finish creating the category object
		category.ID = primitive.NewObjectID()
		category.Category_id = category.ID.Hex()
		category.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		category.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// Insert this category into our category collection in Mongo DB
		result, insertErr := categoryCollection.InsertOne(ctx, category)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, could not create this category."})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var category models.Category
		categoryId := c.Param("category_id")

		// BindJSON reads the body buffer (JSON from Postman) to deserialize it into a Golang struct
		err := c.BindJSON(&category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create an update object and append all necessary details to it
		var updateObj primitive.D

		if category.Name != "" {
			updateObj = append(updateObj, bson.E{"name", category.Name})
		}

		// Set update object's 'Updated_at' field
		category.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", category.Updated_at})

		// Update w MongoDB
		upsert := true
		filter := bson.M{"category_id": categoryId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := categoryCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, failed to update this category."})
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
