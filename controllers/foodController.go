package controller

// getfoods() uses the concept of 'skip and limit'

import (
	"Golang-Management-Platform/database"
	"Golang-Management-Platform/models"
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10 // basically how many records to send per page, pagination basically
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1 // basically we're saying that on the first page we want 10 records (if someone has not sent the page)
		}

		// logic to create pagination
		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		// aggregation in Mongo DB, using pipeline stage functions (match, group, and project)
		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"food_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
				}}}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing food items."})
		}

		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allFoods[0])
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id")
		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing a food item."})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var category models.Category
		var food models.Food

		// BindJSON reads the body buffer (JSON from Postman) to deserialize it into a Golang struct
		err := c.BindJSON(&food)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Handle validation (validate the data based on our user struct)
		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if a category with that category_id already exists
		err = categoryCollection.FindOne(ctx, bson.M{"category_id": food.Category_id}).Decode(&category)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, could not find this category."})
			return
		}

		// Finish creating the food object
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		num := toFixed(*food.Price, 2)
		food.Price = &num

		// Insert this food item into our food collection in Mongo DB
		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, could not create this food item."})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var category models.Category
		var food models.Food

		foodId := c.Param("food_id")

		// BindJSON reads the body buffer (JSON from Postman) to deserialize it into a Golang struct
		err := c.BindJSON(&food)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create an update object and append all necessary details to it
		var updateObj primitive.D

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{"name", food.Name})
		}

		if food.Price != nil {
			updateObj = append(updateObj, bson.E{"price", food.Price})
		}

		if food.Food_image != nil {
			updateObj = append(updateObj, bson.E{"food_image", food.Food_image})
		}

		if food.Category_id != nil {
			err := categoryCollection.FindOne(ctx, bson.M{"category_id": food.Category_id}).Decode(&category)
			defer cancel()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, could not find this category."})
				return
			}
			updateObj = append(updateObj, bson.E{"category_id", food.Category_id})
		}

		// Set update object's 'Updated_at' field
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", food.Updated_at})

		// Update w MongoDB
		upsert := true
		filter := bson.M{"food_id": foodId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := foodCollection.UpdateOne(
			ctx,    // context is basically used to update that particular food item's data
			filter, // filtering using user_id
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, failed to update food item."})
			return
		}

		// All good
		c.JSON(http.StatusOK, result)
	}
}
