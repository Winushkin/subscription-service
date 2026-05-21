// Package entities содержит все сущности приложения
package entities

import (
	"time"
)

// Subscription - запись о подписке пользовтаеля
type Subscription struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServiceName string    `json:"service_name" example:"Yandex Plus"`
	Price       int       `json:"price" example:"400"`
	UserID      string    `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string    `json:"start_date" example:"07-2025"`
	EndDate     *string   `json:"end_date,omitempty" example:"12-2025"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSubscriptionRequest - запрос на создание подписки пользователя
type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price       int     `json:"price" binding:"required,min=1" example:"400"`
	UserID      string  `json:"user_id" binding:"required,uuid" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string  `json:"start_date" binding:"required" example:"07-2025"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2025"`
}

// UpdateSubscriptionRequest - запрос на обновление подписки пользователя
type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty" example:"Yandex Plus"`
	Price       *int    `json:"price,omitempty" example:"500"`
	StartDate   *string `json:"start_date,omitempty" example:"07-2025"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2025"`
}

// CostReportRequest - запрос на отчет по суммарной стоимости подписок пользователя
type CostReportRequest struct {
	UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string `json:"service_name" example:"Yandex Plus"`
	StartDate   string `json:"start_date" binding:"required" example:"06-2025"`
	EndDate     string `json:"end_date" binding:"required" example:"08-2025"`
}

// CostReport - отчет по суммарной стоимости подписок пользователя
type CostReport struct {
	TotalCost int    `json:"total_cost" example:"1200"`
	Count     int    `json:"count" example:"3"`
	Currency  string `json:"currency" example:"RUB"`
}

// ErrorResponse - ответ возвращаемой ошибки
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request body"`
}