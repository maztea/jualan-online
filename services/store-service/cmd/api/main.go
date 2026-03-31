package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"jualan-online/services/common/config"
	"jualan-online/services/common/database"
	"jualan-online/services/common/logger"
	"jualan-online/services/store-service/internal/domain"
	"jualan-online/services/store-service/internal/handler"
	"jualan-online/services/store-service/internal/repository"
	"jualan-online/services/store-service/internal/service"
	"os"
	"os/signal"
	"time"
)

func main() {
	// 1. Init Logger
	logger.InitLogger()
	log := logger.GetLogger()

	// 2. Database connection
	dbConfig := database.Config{
		Host:     config.GetEnv("DB_HOST", "localhost"),
		Port:     config.GetEnv("DB_PORT", "5432"),
		User:     config.GetEnv("DB_USER", "postgres"),
		Password: config.GetEnv("DB_PASSWORD", "postgres"),
		DBName:   config.GetEnv("DB_NAME", "store_db"),
		SSLMode:  config.GetEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.Connect(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}

	// Auto migration
	if err := db.AutoMigrate(&domain.Store{}); err != nil {
		logger.Fatal("Failed to migrate database: " + err.Error())
	}

	// 3. Init Layers
	repo := repository.NewStoreRepository(db)
	svc := service.NewStoreService(repo)
	hdl := handler.NewStoreHandler(svc)

	// 4. Setup Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	v1 := e.Group("/v1")
	{
		v1.GET("/stores", hdl.GetAll)
		v1.POST("/stores", hdl.Create)
		v1.GET("/stores/:id", hdl.GetByID)
		v1.PUT("/stores/:id", hdl.Update)
		v1.DELETE("/stores/:id", hdl.Delete)
	}
	e.GET("/health", hdl.Health)

	// 5. Start Server
	go func() {
		if err := e.Start(":" + config.GetEnv("PORT", "8082")); err != nil {
			log.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
