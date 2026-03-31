package repository

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"jualan-online/services/auth-service/internal/domain"
	"jualan-online/services/common/config"
	"jualan-online/services/common/logger"
)

func SeedUsers(db *gorm.DB) error {
	logger.Info("Starting to seed users...")

	// Check if users already exist
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 0 {
		logger.Info("Users already seeded, skipping...")
		return nil
	}

	adminPasswordRaw := config.GetEnv("SEED_ADMIN_PASSWORD", "admin123")
	staffPasswordRaw := config.GetEnv("SEED_STAFF_PASSWORD", "staff123")

	adminPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPasswordRaw), bcrypt.DefaultCost)
	staffPassword, _ := bcrypt.GenerateFromPassword([]byte(staffPasswordRaw), bcrypt.DefaultCost)

	users := []domain.User{
		{
			Username:     "admin",
			PasswordHash: string(adminPassword),
			Role:         "ADMIN",
		},
		{
			Username:     "staff",
			PasswordHash: string(staffPassword),
			Role:         "STAFF",
		},
	}

	for _, u := range users {
		if err := db.Create(&u).Error; err != nil {
			return err
		}
		logger.Info("Seeded user: " + u.Username)
	}

	logger.Info("Seeding completed successfully")
	return nil
}
