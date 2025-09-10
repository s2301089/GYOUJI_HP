// @title GYOUJI_HP API
// @version 1.0
// @description This is a sample server for GYOUJI_HP.
// @host localhost:8080
// @BasePath /
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"
	"github.com/saku0512/GYOUJI_HP/backend/internal/router"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("Starting the application...")
	// Database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_ROOT_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database")

	// Initialize users
	initializeUsers(db)

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY is not set")
	}

	// Repository -> Service -> Handler の順でインスタンスを生成
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, jwtSecret)
	userHandler := handler.NewUserHandler(userService)

	r := router.SetupRouter(userHandler, jwtSecret)
	log.Println("Server is running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func initializeUsers(db *sql.DB) {
	users := []map[string]string{
		{"username": os.Getenv("SUPERROOT_USER"), "password": os.Getenv("SUPERROOT_PASS"), "role": "superroot", "assigned_sport": ""},
		{"username": os.Getenv("ADMIN_TABLE_TENNIS_USER"), "password": os.Getenv("ADMIN_TABLE_TENNIS_PASS"), "role": "admin", "assigned_sport": "table_tennis"},
		{"username": os.Getenv("ADMIN_VOLLEYBALL_USER"), "password": os.Getenv("ADMIN_VOLLEYBALL_PASS"), "role": "admin", "assigned_sport": "volleyball"},
		{"username": os.Getenv("ADMIN_SOCCER_USER"), "password": os.Getenv("ADMIN_SOCCER_PASS"), "role": "admin", "assigned_sport": "soccer"},
		{"username": os.Getenv("STUDENT_USER"), "password": os.Getenv("STUDENT_PASS"), "role": "student", "assigned_sport": ""},
	}

	for _, user := range users {
		log.Printf("Processing user: %s", user["username"])
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", user["username"]).Scan(&exists)
		if err != nil {
			log.Printf("failed to check if user %s exists: %v", user["username"], err)
			continue
		}

		log.Printf("User %s exists: %t", user["username"], exists)

		if !exists {
			log.Printf("User %s does not exist, creating...", user["username"])
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user["password"]), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("failed to hash password for user %s: %v", user["username"], err)
				continue
			}

			_, err = db.Exec("INSERT INTO users (username, hashed_password, role, assigned_sport) VALUES (?, ?, ?, ?)",
				user["username"], string(hashedPassword), user["role"], user["assigned_sport"])
			if err != nil {
				log.Printf("failed to insert user %s: %v", user["username"], err)
				continue
			}
			fmt.Printf("User %s created successfully.\n", user["username"])
		}
	}
}
