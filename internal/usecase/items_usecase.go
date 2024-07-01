package usecase

import (
	"context"
	"fmt"
	"github.com/robertgarayshin/warehousesAPI/internal/entity"
)

// ItemsUseCase -.
type ItemsUseCase struct {
	itemsRepository ItemsRepo
}

// NewItemsUsecase -.
func NewItemsUsecase(i ItemsRepo) ItemsUseCase {
	return ItemsUseCase{
		itemsRepository: i,
	}
}

// CreateItems -
func (uc *ItemsUseCase) CreateItems(ctx context.Context, items []entity.Item) error {
	if err := uc.itemsRepository.StoreItems(ctx, items); err != nil {
		return fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}

	return nil
}

func (uc *ItemsUseCase) Quantity(ctx context.Context, id int) (map[string]int, error) {
	quantity, err := uc.itemsRepository.QuantityByWarehouse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error get quantity by warehouse %d: %w", id, err)
	}

	return quantity, nil
}
