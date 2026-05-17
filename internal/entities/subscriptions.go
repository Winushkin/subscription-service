// Package entities хранит сущности приложения
package entities

import (
	"time"
)

type Subscription struct {
	ID          string    `json:"id" db:"id"`
	ServiceName string    `json:"service_name" db:"service_name"`
	Price       int       `json:"price" db:"price"`
	UserID      string    `json:"user_id" db:"user_id"`
	StartDate   string    `json:"start_date" db:"start_date"`       // MM-YYYY
	EndDate     *string   `json:"end_date,omitempty" db:"end_date"` // NULL-able
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int     `json:"price" binding:"required,min=1"`
	UserID      string  `json:"user_id" binding:"required,uuid"`
	StartDate   string  `json:"start_date" binding:"required"` // MM-YYYY
	EndDate     *string `json:"end_date"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name"`
	Price       *int    `json:"price"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

type CostReportRequest struct {
	UserID      string `json:"user_id" binding:"uuid"`
	ServiceName string `json:"service_name"`
	StartDate   string `json:"start_date" binding:"required"` // MM-YYYY
	EndDate     string `json:"end_date" binding:"required"`   // MM-YYYY
}

type CostReport struct {
	TotalCost int    `json:"total_cost"`
	Count     int    `json:"count"`
	Currency  string `json:"currency"`
}
