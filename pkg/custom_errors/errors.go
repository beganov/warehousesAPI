package custom_errors

import "errors"

var (
	ErrWarehouseUnavailable = errors.New("warehouse is unavailable")
	ErrNoReservation        = errors.New("no reservation found")
)
