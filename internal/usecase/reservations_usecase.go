package usecase

import (
	"context"
	"fmt"
)

type ReservationsUsecase struct {
	reservationsRepository ReservationsRepo
}

func (uc *ReservationsUsecase) Reserve(ctx context.Context, ids []string) error {
	if err := uc.reservationsRepository.CreateReservation(ctx, ids); err != nil {
		return fmt.Errorf("error create reservation. %w", err)
	}

	return nil
}

func (uc *ReservationsUsecase) CancelReservation(ctx context.Context, ids []string) error {
	if err := uc.reservationsRepository.DeleteReservation(ctx, ids); err != nil {
		return fmt.Errorf("error delete reservation. %w", err)
	}

	return nil
}

func NewReservationsUsecase(reservationsRepository ReservationsRepo) ReservationsUsecase {
	return ReservationsUsecase{
		reservationsRepository: reservationsRepository,
	}
}
