package repo

import (
	"context"
	"fmt"
	"github.com/robertgarayshin/warehousesAPI/internal/entity"
	"github.com/robertgarayshin/warehousesAPI/pkg/postgres"
)

const _defaultEntityCap = 64

// ItemsRepo -.
type ItemsRepo struct {
	*postgres.Postgres
}

func (r *ItemsRepo) StoreItems(ctx context.Context, items []entity.Item) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction. %w", err)
	}

	for _, item := range items {
		stmt := `INSERT INTO items (unique_code, name, size, quantity, warehouse_id) 
							VALUES ($1, $2, $3, $4, $5)
								ON CONFLICT (unique_code) DO UPDATE
								SET name = $2, 
								    size = $3, 
								    quantity = items.quantity + $4, 
								    warehouse_id = $5
								WHERE items.unique_code = $1`
		_, err = tx.Exec(ctx, stmt, item.UniqueId, item.Name, item.Size, item.Quantity, item.WarehouseId)
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return fmt.Errorf("transaction already closed. %w", err)
			}

			return fmt.Errorf("error executing insert item statement. %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf("transaction already closed. %w", err)
		}
		return fmt.Errorf("error commiting transaction. %w", err)
	}

	return nil
}

func (r *ItemsRepo) QuantityByWarehouse(ctx context.Context, warehouseId int) (map[string]int, error) {
	res := make(map[string]int)
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	stmt, args, err := r.Builder.Select("unique_code, quantity").
		From("items").
		Where("warehouse_id = ?", warehouseId).ToSql()
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return nil, fmt.Errorf("transaction already closed. %w", err)
		}

		return nil, fmt.Errorf("error building query. %w", err)
	}

	rows, err := tx.Query(ctx, stmt, args...)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return nil, fmt.Errorf("transaction already closed. %w", err)
		}

		return nil, fmt.Errorf("error executing query. %w", err)
	}

	for rows.Next() {
		var id string
		var quantity int

		err := rows.Scan(&id, &quantity)
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return nil, fmt.Errorf("transaction already closed. %w", err)
			}

			return nil, fmt.Errorf("error scanning row value. %w", err)
		}

		res[id] = quantity
	}

	if err := tx.Commit(ctx); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, fmt.Errorf("transaction already closed. %w", err)
		}

		return nil, fmt.Errorf("error commiting transaction. %w", err)
	}

	return res, nil
}

// NewItemsRepository -.
func NewItemsRepository(pg *postgres.Postgres) *ItemsRepo {
	return &ItemsRepo{pg}
}
