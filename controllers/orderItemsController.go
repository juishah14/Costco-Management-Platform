package controller

// GET RID OF THE ORDER ITEM PACK, just have this be regular

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

type OrderItemPack struct {
	Membership_id *string
	Order_items   []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

// Will implement this func later
func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	return OrderItems, err
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")
		allOrderItems, err := ItemsByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing order items by their order ID."})
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := orderItemCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing order items."})
			return
		}

		var allOrderItems []bson.M
		err = result.All(ctx, &allOrderItems)
		if err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItem models.OrderItem
		orderItemId := c.Param("order_item_id")

		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing this order item."})
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var order models.Order
		var orderItemPack OrderItemPack

		// BindJSON reads the body buffer (JSON from Postman) to deserialize it into a Golang struct
		err := c.BindJSON(&orderItemPack)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Finish creating the order object
		order_id := OrderItemOrderCreator(order) // Create a new order for this order item
		order.Membership_id = orderItemPack.Membership_id
		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItemsToBeInserted := []interface{}{}

		for _, orderItem := range orderItemPack.Order_items {
			orderItem.ID = primitive.NewObjectID()
			orderItem.Order_item_id = orderItem.ID.Hex()
			orderItem.Order_id = order_id
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

			// Handle validation (validate the data based on our orderItem struct)
			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			// Add order item to our slice
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		// Insert all order items in our slice to our order item collection in Mongo DB
		insertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			log.Fatal(err)
		}
		defer cancel()
		c.JSON(http.StatusOK, insertedOrderItems)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItem models.OrderItem
		orderItemId := c.Param("order_item_id")

		// Create an update object and append all necessary details to it
		var updateObj primitive.D

		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity", *orderItem.Quantity})
		}

		if orderItem.Product_id != nil {
			updateObj = append(updateObj, bson.E{"product_id", *orderItem.Product_id})
		}

		if orderItem.Order_id != "" {
			updateObj = append(updateObj, bson.E{"order_id", orderItem.Order_id})
		}

		// Set update object's 'Updated_at' field
		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", orderItem.Updated_at})

		// Update w Mongo DB
		upsert := true
		filter := bson.M{"order_item_id": orderItemId}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error, failed to update this order item."})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
