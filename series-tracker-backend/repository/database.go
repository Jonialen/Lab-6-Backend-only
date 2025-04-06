// Package repository maneja toda la lógica de interacción con la base de datos,
// abstraendo las operaciones de GORM del resto de la aplicación.
package repository

import (
	"fmt"
	"log"
	"os" // Importar os para leer variables de entorno

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lab6/models" // Asegúrate que la ruta de importación sea correcta
)

// DB es la instancia global y exportada de la conexión a la base de datos GORM.
// Otros paquetes la usarán para realizar operaciones en la base de datos.
var DB *gorm.DB

// InitDB inicializa la conexión con la base de datos MySQL usando GORM.
// Lee la cadena de conexión (DSN) preferentemente desde variables de entorno
// y realiza la automigración del modelo Series.
// Termina la aplicación si la conexión o la migración fallan.
func InitDB() {
	// Leer configuración de la base de datos desde variables de entorno
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Valores por defecto si no se establecen las variables de entorno
	// (Útil para desarrollo local si no se usa Docker Compose con env vars)
	if dbUser == "" {
		dbUser = "app_user"
	}
	if dbPassword == "" {
		dbPassword = "app_password" // ¡Cambiar en producción!
	}
	if dbHost == "" {
		dbHost = "database" // Nombre del servicio en Docker Compose o localhost
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbName == "" {
		dbName = "anime_db"
	}

	// Construir la cadena de conexión (DSN)
	// Añadidos parámetros recomendados: charset, parseTime, loc
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error fatal al conectar a la base de datos: %v\nDSN: %s", err, dsn)
	}

	fmt.Println("Conexión a la base de datos exitosa.")

	// AutoMigrate intentará crear o actualizar la tabla 'series' según el modelo models.Series.
	fmt.Println("Ejecutando AutoMigrate para el modelo Series...")
	err = DB.AutoMigrate(&models.Series{})
	if err != nil {
		log.Fatalf("Error fatal durante AutoMigrate: %v", err)
	}
	fmt.Println("AutoMigrate completado.")
}

// CloseDB cierra la conexión a la base de datos si está abierta.
// Es importante llamar a esta función para liberar recursos idealmente durante el cierre grácil de la aplicación.
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Error obteniendo la instancia DB subyacente: %v", err)
		}
		if sqlDB != nil {
			fmt.Println("Cerrando conexión a la base de datos...")
			if err := sqlDB.Close(); err != nil {
				log.Printf("Error cerrando la conexión a la base de datos: %v", err)
			} else {
				fmt.Println("Conexión a la base de datos cerrada.")
			}
		}
	}
}
