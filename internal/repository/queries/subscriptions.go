// Package queries содержит функции для сборки запросов к базе данных
package queries

import (
	"fmt"
	"subscription-service/internal/entities"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const (
	// Константы с именами таблиц и колонок используемых в запросах
	subTable    = "subscriptions"
	subID       = "id"
	serviceName = "service_name"
	price       = "price"
	userID      = "user_id"
	startDate   = "start_date"
	endDate     = "end_date"
	createdAt   = "created_at"
	updatedAt   = "updated_at"
)

// psql - это экземпляр StatementBuilder, настроенный на использование формата плейсхолдеров для PostgreSQL.
var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// Exists возвращает pапрос на проверку наличия сущности пива по id
// func Exists(id uint) sq.SelectBuilder {
// 	return psql.Select("id").From(beersTable).Where(sq.Eq{"id": id})
// }

func columnList() string {
	return fmt.Sprintf(
		"%s, %s, %s, %s, %s, %s, %s",
		subID,
		serviceName,
		price,
		startDate,
		endDate,
		createdAt,
		updatedAt,
	)
}

// FullSubSelect возвращает базовый запрос для получения полной информации о пиве, включая его характеристики, город и страну производства, категорию и особенности.
func FullSubSelect() sq.SelectBuilder {
	return psql.Select(
		subID,
		serviceName,
		price,
		startDate,
		endDate,
		createdAt,
		updatedAt,
	).From(subTable)
}

// InsertSubscription возвращает запрос для вставки нового отзыва в таблицу reviews и возвращает ID вставленного отзыва.
func InsertSubscription(sub entities.CreateSubscriptionRequest) sq.InsertBuilder {
	data := map[string]any{
		serviceName: sub.ServiceName,
		price:       sub.Price,
		userID:      sub.UserID,
		startDate:   sub.StartDate,
		endDate:     sub.EndDate,
	}

	return psql.
		Insert(subTable).
		SetMap(data).
		Suffix("RETURNING " + columnList())
}

func SelectSubByID(id uuid.UUID) sq.SelectBuilder {
	return FullSubSelect().Where(sq.Eq{subID: id})
}

func UpdateSubscription(id uuid.UUID, req entities.UpdateSubscriptionRequest) sq.UpdateBuilder {
	updates := make(map[string]any)

	if req.ServiceName != nil {
		updates[serviceName] = req.ServiceName
	}
	if req.Price != nil {
		updates[price] = req.Price
	}
	if req.StartDate != nil {
		updates[startDate] = req.StartDate
	}
	if req.EndDate != nil {
		updates[endDate] = req.EndDate
	}

	return psql.Update(subTable).
		SetMap(updates).
		Where(sq.Eq{subID: id})
}

func DeleteSubscription(id uuid.UUID) sq.DeleteBuilder {
	return psql.
		Delete(subTable).
		Where(sq.Eq{subID: id})
}


func SelectSubscriptionsCost(req entities.CostReportRequest) sq.SelectBuilder{
	sumSelection := fmt.Sprintf("SUM(%s)", price)
	query := psql.Select(sumSelection).From(subTable)

	if req.ServiceName != ""{
		query = query.Where(sq.Eq{serviceName: req.ServiceName})
	}
	if req.StartDate != ""{
		query = query.Where(sq.LtOrEq{startDate: req.StartDate})
	}
	if req.EndDate != ""{
		query = query.Where(sq.GtOrEq{endDate: req.EndDate})
	}
	
	return query
}