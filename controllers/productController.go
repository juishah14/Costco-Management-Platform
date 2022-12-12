package controller

// used to be food

import (
	"Golang-Management-Platform/database"
	"Golang-Management-Platform/models"
	"context"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var productCollection *mongo.Collection = database.OpenCollection(database.Client, "product")
var validate = validator.New()

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := productCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing products."})
		}

		var allProducts []bson.M
		err = result.All(ctx, &allProducts)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allProducts)
	}
}

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var product models.Product
		productId := c.Param("product_id")

		err := productCollection.FindOne(ctx, bson.M{"product_id": productId}).Decode(&product)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing a product."})
		}
		c.JSON(http.StatusOK, product)
	}
}

func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var product models.Product
		var category models.Category

		// BindJSON reads the body buffer (JSON from Postman) to deserialize it into a Golang struct
		err := c.BindJSON(&product)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Handle validation (validate the data based on our product struct)
		validationErr := validate.Struct(product)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if a category with that category name already exists
		err = categoryCollection.FindOne(ctx, bson.M{"name": product.Category}).Decode(&category)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, could not find this category."})
			return
		}

		// Finish creating the product object
		product.ID = primitive.NewObjectID()
		product.Product_id = product.ID.Hex()
		num := toFixed(*product.Price, 2)
		product.Price = &num
		product.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// Insert this product into our product collection in Mongo DB
		result, insertErr := productCollection.InsertOne(ctx, product)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, could not create this product."})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var product models.Product
		var category models.Category
		productId := c.Param("product_id")

		// BindJSON reads the body buffer (JSON from Postman) to deserialize it into a Golang struct
		err := c.BindJSON(&product)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create an update object and append all necessary details to it
		var updateObj primitive.D

		if product.Name != nil {
			updateObj = append(updateObj, bson.E{"name", product.Name})
		}

		if product.Price != nil {
			updateObj = append(updateObj, bson.E{"price", product.Price})
		}

		if product.Description != nil {
			updateObj = append(updateObj, bson.E{"description", product.Description})
		}

		if product.Category != nil {
			err := categoryCollection.FindOne(ctx, bson.M{"name": product.Category}).Decode(&category)
			defer cancel()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, could not find this category."})
				return
			}
			updateObj = append(updateObj, bson.E{"category", product.Category})
		}

		// Set update object's 'Updated_at' field
		product.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", product.Updated_at})

		// Update w MongoDB
		upsert := true
		filter := bson.M{"product_id": productId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := productCollection.UpdateOne(
			ctx,    // context is basically used to update that particular product's data
			filter, // filtering using user_id
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, failed to update this product."})
			return
		}

		// All good
		c.JSON(http.StatusOK, result)
	}
}
