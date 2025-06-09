package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Loan struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VehicleID        primitive.ObjectID `json:"vehicle_id" bson:"vehicle_id"`
	CompanyID        primitive.ObjectID `json:"company_id" bson:"company_id"`
	Lender           string             `json:"lender" bson:"lender" validate:"required"`
	PrincipalAmount  float64            `json:"principal_amount" bson:"principal_amount" validate:"required,min=0"`
	InterestRate     float64            `json:"interest_rate" bson:"interest_rate" validate:"required,min=0,max=100"`
	TermMonths       int                `json:"term_months" bson:"term_months" validate:"required,min=1"`
	StartDate        time.Time          `json:"start_date" bson:"start_date"`
	MonthlyPayment   float64            `json:"monthly_payment" bson:"monthly_payment"`
	RemainingBalance float64            `json:"remaining_balance" bson:"remaining_balance"`
	Status           string             `json:"status" bson:"status" validate:"required,oneof=active paid_off"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

type Payment struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	LoanID           primitive.ObjectID `json:"loan_id" bson:"loan_id"`
	PaymentDate      time.Time          `json:"payment_date" bson:"payment_date"`
	PrincipalPaid    float64            `json:"principal_paid" bson:"principal_paid"`
	InterestPaid     float64            `json:"interest_paid" bson:"interest_paid"`
	TotalPaid        float64            `json:"total_paid" bson:"total_paid"`
	RemainingBalance float64            `json:"remaining_balance" bson:"remaining_balance"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
}
