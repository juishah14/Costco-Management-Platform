package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// used to be food

type Product struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        *string            `json:"name" validate:"required,min=2,max=100"`
	Price       *float64           `json:"price" validate:"required"`
	Description *string            `json:"description" validate:"required"` // instead of food_image
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	Product_id  string             `json:"product_id"`
	Category    *string            `json:"category" validate:"required"` // used to be menu_id, category name instead of id
}
