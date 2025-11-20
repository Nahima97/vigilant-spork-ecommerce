package db

import (
    "fmt"
    "log"
	"vigilant-spork/models"
    "os"
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() *gorm.DB {

    var err error

    err = godotenv.Load()
    if err != nil {
        log.Fatalf("error loading .env file: %v", err)
    }

    connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    Db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    fmt.Println("connected to database successfully!")

    err = Db.AutoMigrate(&models.User{}, &models.Product{}, &models.Cart{}, &models.CartItem{}, &models.Order{}, &models.OrderItem{}, &models.Review{})
    if err != nil {
        log.Fatalf("unable to migrate schema: %v", err)
    }
    fmt.Println("Database automigration completed!")
    return Db
}