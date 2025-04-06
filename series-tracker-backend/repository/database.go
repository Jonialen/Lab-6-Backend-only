package repository

import (
	"fmt"
	"log"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lab6/models"
)

var DB *gorm.DB

// InitDB inicializa la conexión con la base de datos
func InitDB() {
	// Modify this to use the database service defined in your docker-compose.yml
	dsn := "app_user:app_password@tcp(database:3306)/anime_db"
	
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	
	// Auto migrate your models
	DB.AutoMigrate(&models.Series{})
	
	fmt.Println("Conexión a la base de datos exitosa")
}

// CloseDB cierra la conexión con la base de datos
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Fatal("Error getting DB instance:", err)
		}
		sqlDB.Close()
	}
}
