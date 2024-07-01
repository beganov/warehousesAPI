package repo

import (
	"context"
	"fmt"
	"warehousesAPI/pkg/postgres"
)

type ReservationsRepo struct {
	*postgres.Postgres
}

func (r *ReservationsRepo) CreateReservation(ctx context.Context, ids []string) error {
	reservations := r.reservedItemsCount(ids)

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error begining transaction. %w", err)
	}

	for id, count := range reservations {
		// Запись склада Unstored (id = 0) для товаров, которые создаются при создании резервации
		/*
			Usecase: товар добавлен в предварительный резерв, но при этом еще не прибыл на склад.
			Записи об этом товаре нет в таблице товаров.
			Создается резервация -> уникальный код товара создает дефолтный товар (все поля пустые, кроме кода)
			Товар помещается на склад unstored.
			Количество товара на складе = 0 - количество товара в резервации (уменьшается при следующей резервации)
			Товар прибывает на склад: запись о нем обновляется, указывается реальный склад, добавляется вся информация
				количество товара = текущее количество + количество прибывшего
			Резервация товара сохраняется, новое количество товара учитывает зарезервированный
		*/
		itemCreateStatement := `INSERT INTO items(unique_code) VALUES ($1) ON CONFLICT DO NOTHING`

		_, err := tx.Exec(ctx, itemCreateStatement, id)
		if err != nil {
			return fmt.Errorf("error create item. %w", err)
		}

		reservationCreateStatement := `UPDATE items 
											SET reserved = reserved + $1,
											    quantity = quantity - $1
										WHERE unique_code = $2`

		_, err = tx.Exec(ctx, reservationCreateStatement, count, id)
		if err != nil {
			return fmt.Errorf("error executing insert reservation. %w", err)
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

func (r *ReservationsRepo) DeleteReservation(ctx context.Context, ids []string) error {
	deleteReservations := r.reservedItemsCount(ids)

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error begining transaction. %w", err)
	}

	for id, count := range deleteReservations {
		reservationDeleteStatement := `UPDATE items
											SET reserved = reserved - $1,
											    quantity = quantity + $1
										WHERE unique_code = $2`

		_, err = tx.Exec(ctx, reservationDeleteStatement, count, id)
		if err != nil {
			// todo: кастомная ошибка
			return fmt.Errorf("error delete reservation. %w", err)
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

func (r *ReservationsRepo) reservedItemsCount(ids []string) map[string]int {
	res := make(map[string]int, len(ids))

	for _, id := range ids {
		res[id] += 1
	}

	return res
}

func NewReservationRepo(p *postgres.Postgres) *ReservationsRepo {
	return &ReservationsRepo{p}
}
