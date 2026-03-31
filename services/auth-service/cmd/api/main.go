package main

import (
	"context"
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"jualan-online/services/auth-service/internal/domain"
	"jualan-online/services/auth-service/internal/handler"
	"jualan-online/services/auth-service/internal/repository"
	"jualan-online/services/auth-service/internal/service"
	"jualan-online/services/common/config"
	"jualan-online/services/common/database"
	"jualan-online/services/common/logger"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Parse Flags
	seedFlag := flag.Bool("seed", false, "run seeding")
	flag.Parse()

	// 1. Init Logger
	logger.InitLogger()
	log := logger.GetLogger()

	// 2. Database connection
	dbConfig := database.Config{
		Host:     config.GetEnv("DB_HOST", "localhost"),
		Port:     config.GetEnv("DB_PORT", "5432"),
		User:     config.GetEnv("DB_USER", "postgres"),
		Password: config.GetEnv("DB_PASSWORD", "postgres"),
		DBName:   config.GetEnv("DB_NAME", "auth_db"),
		SSLMode:  config.GetEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.Connect(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}

	// Auto migration
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		logger.Fatal("Failed to migrate database: " + err.Error())
	}

	// 3. Check Seed Flag
	if *seedFlag {
		if err := repository.SeedUsers(db); err != nil {
			logger.Fatal("Failed to seed users: " + err.Error())
		}
		os.Exit(0)
	}

	// 4. Init Layers
	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(repo, config.GetEnv("JWT_SECRET", "very-secret-key"))
	hdl := handler.NewAuthHandler(svc)

	// 4. Setup Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health", hdl.Health)
	e.POST("/register", hdl.Register)
	e.POST("/login", hdl.Login)

	// 5. Start Server
	go func() {
		if err := e.Start(":" + config.GetEnv("PORT", "8081")); err != nil {
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
