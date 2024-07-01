package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
	"warehousesAPI/config"
	v1 "warehousesAPI/internal/controller/http/v1"
	"warehousesAPI/internal/usecase"
	"warehousesAPI/internal/usecase/repo"
	"warehousesAPI/pkg/httpserver"
	"warehousesAPI/pkg/logger"
	"warehousesAPI/pkg/postgres"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.NewItemsRepository: %w", err))
	}
	defer pg.Close()

	// Use case
	itemsUsecase := usecase.NewItemsUsecase(
		repo.NewItemsRepository(pg),
	)

	reservationsUsecase := usecase.NewReservationsUsecase(
		repo.NewReservationRepo(pg),
	)

	warehousesUsecase := usecase.NewWarehousesUsecase(
		repo.NewWarehousesRepo(pg),
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler,
		l,
		itemsUsecase,
		reservationsUsecase,
		warehousesUsecase,
	)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
