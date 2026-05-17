// Package repository содержит слой для манипуляции объектами в базе данных
package repository

import (
	"context"
	"errors"
	"fmt"
	"subscription-service/internal/entities"
	"subscription-service/internal/repository/queries"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	// CreateSubscription создает новую запись подписки в БД и возвращает его
	CreateSubscription(ctx context.Context, req entities.CreateSubscriptionRequest) (*entities.Subscription, error)

	// GetSubscription возвращает объект entities.Subscription по id
	GetSubscription(ctx context.Context, id uuid.UUID) (*entities.Subscription, error)

	// GetSubscription возвращает слайс объектов entities.Subscription по id пользователя
	ListSubscriptions(ctx context.Context, userID uuid.UUID, limit, offset uint64) ([]entities.Subscription, error)

	// UpdateSubscription обновляет запись о подписке в БД
	UpdateSubscription(ctx context.Context, id uuid.UUID, req entities.UpdateSubscriptionRequest) (*entities.Subscription, error)

	// GetSubscriptionsCost возвращает отчет о стоимости подписок по фильтрам
	GetSubscriptionsCost(ctx context.Context, req entities.CostReportRequest) (*entities.CostReport, error)
}

const currency = "RUB"

type Postgres struct {
	pool *pgxpool.Pool
}

// NewRepository создает новый экземпляр Repository с переданным пулом соединений.
func NewRepository(pool *pgxpool.Pool) Repository {
	return &Postgres{pool: pool}
}

func (p *Postgres) CreateSubscription(ctx context.Context, req entities.CreateSubscriptionRequest) (*entities.Subscription, error) {
	psql := queries.InsertSubscription(req)
	query, args, err := psql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	row := p.pool.QueryRow(ctx, query, args...)
	createdSub, err := scanSub(row)

	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w ", err)
	}
	return createdSub, nil
}

func (p *Postgres) GetSubscription(ctx context.Context, id uuid.UUID) (*entities.Subscription, error) {
	psql := queries.SelectSubByID(id)
	query, args, err := psql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	row := p.pool.QueryRow(ctx, query, args...)
	sub, err := scanSub(row)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	return sub, nil
}

func (p *Postgres) ListSubscriptions(ctx context.Context, userID uuid.UUID, limit, offset uint64) ([]entities.Subscription, error) {
	psql := queries.FullSubSelect().Offset(offset)
	if limit != 0 {
		psql = psql.Limit(limit)
	}

	query, args, err := psql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	subs := make([]entities.Subscription, 0)
	for rows.Next() {
		sub, err := scanSub(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		subs = append(subs, *sub)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return subs, nil
}

func (p *Postgres) UpdateSubscription(ctx context.Context, id uuid.UUID, req entities.UpdateSubscriptionRequest) (*entities.Subscription, error) {
	psql := queries.UpdateSubscription(id, req)
	query, args, err := psql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	row := p.pool.QueryRow(ctx, query, args...)
	createdSub, err := scanSub(row)

	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w ", err)
	}
	return createdSub, nil
}

func (p *Postgres) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	updatePsql := queries.DeleteSubscription(id)
	query, args, err := updatePsql.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", "toSql", err)
	}

	result, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", "exec", err)
	}

	if result.RowsAffected() != 1 {
		return errors.New("failed to delete subscription")
	}

	return nil
}

func (p *Postgres) GetSubscriptionsCost(ctx context.Context, req entities.CostReportRequest) (*entities.CostReport, error) {
	updatePsql := queries.SelectSubscriptionsCost(req)
	query, args, err := updatePsql.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "toSql", err)
	}

	var cost int
	err = p.pool.QueryRow(ctx, query, args...).Scan(&cost)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "exec", err)
	}

	report := &entities.CostReport{
		TotalCost: cost,
		Currency:  currency,
	}

	return report, nil
}

// scanSub сканирует полную сырую строку из базы данных в сущность entities.Subscription
func scanSub(row pgx.Row) (*entities.Subscription, error) {
	var sub entities.Subscription
	err := row.Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "Scan", err)
	}
	return &sub, nil
}
