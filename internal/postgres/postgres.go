package postgres

import (
	"order_svc/config"
	"order_svc/internal/models"

	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	DB *gorm.DB
}

func (postgresDB *PostgresDB) Init() {
	// using gorm to connect to db driver
	configs := config.ReadConfigs(".")
	var dns string

	databaseUrl := configs.DbUrl
	// os.Getenv("DATABASE_URL")

	if databaseUrl == "" {
		dns = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			configs.DbHost, configs.DbPort,
			configs.DbUser, configs.DbPassword, configs.DbName)
	} else {
		dns = databaseUrl
	}
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	postgresDB.DB = db

	//err = db.Migrator().DropTable(&models.CartItem{}, &models.CartPackageItem{}, &models.NewOrder{})
	err = db.AutoMigrate(&models.NewOrder{}, &models.CartItem{}, &models.CartPackageItem{})
	if err != nil {
		log.Panic(err)
	}

	log.Println("Database Connected Successfully...")
}
