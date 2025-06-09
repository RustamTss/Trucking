package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vehicle struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CompanyID     primitive.ObjectID `json:"company_id" bson:"company_id"`
	Type          string             `json:"type" bson:"type" validate:"required,oneof=truck trailer"`
	VIN           string             `json:"vin" bson:"vin" validate:"required"`
	Make          string             `json:"make" bson:"make" validate:"required"`
	Model         string             `json:"model" bson:"model" validate:"required"`
	Year          int                `json:"year" bson:"year" validate:"required,min=1900,max=2030"`
	PurchasePrice float64            `json:"purchase_price" bson:"purchase_price" validate:"required,min=0"`
	PurchaseDate  time.Time          `json:"purchase_date" bson:"purchase_date"`
	Status        string             `json:"status" bson:"status" validate:"required,oneof=active inactive sold"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}
