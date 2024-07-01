package usecase

import (
	"context"
	"fmt"
	"github.com/robertgarayshin/warehousesAPI/internal/entity"
)

type WarehousesUsecase struct {
	warehousesRepository WarehousesRepo
}

func (w *WarehousesUsecase) WarehouseCreate(ctx context.Context, warehouse entity.Warehouse) error {
	if err := w.warehousesRepository.CreateWarehouse(ctx, warehouse); err != nil {
		return fmt.Errorf("error creating warehouse: %w", err)
	}

	return nil
}

// NewWarehousesUsecase -.
func NewWarehousesUsecase(w WarehousesRepo) WarehousesUsecase {
	return WarehousesUsecase{
		warehousesRepository: w,
	}
}
