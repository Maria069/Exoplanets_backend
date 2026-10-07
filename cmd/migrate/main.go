package main

import (
	"exoplanets/internal/app/ds"
	"exoplanets/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.Exoplanet{},
		&ds.User{},
		&ds.Likes{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
