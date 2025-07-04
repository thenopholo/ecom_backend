package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thenopholo/ecom_backend/cmd/api"
	"github.com/thenopholo/ecom_backend/internal/config"
	"github.com/thenopholo/ecom_backend/internal/database"
)

func main() {
	log.Println("Init application backend")
	log.Printf("Env: %s", getEnvironment())

	dbConfig := database.PostgresSQLConfig{
		Host:     config.Envs.DBHost,
		Port:     config.Envs.DBPort,
		User:     config.Envs.DBUser,
		Password: config.Envs.DBPassword,
		DBName:   config.Envs.DBName,
		SSLMode:  config.Envs.SSLMode,
	}

	db, err := connectWithRetry(dbConfig, 5, 2*time.Second)
	if err != nil {
		log.Fatalf("Fatal error on trying to conect to PostgresSQL many times: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error on closing DB conection: %v", err)
		} else {
			log.Println("DB conection closed")
		}
	}()

	if err := initDatabase(db); err != nil {
		log.Fatalf("error on initializing DB: %v", err)
	}

	server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), db)

	go shutdown(server)

	log.Printf("Server running on %s:%s", config.Envs.PublicHost, config.Envs.Port)
	if err := server.Run(); err != nil {
		log.Fatalf("Error on initializing server: %v", err)
	}
}

// func initStorage(db *sql.DB) {
// 	err := db.Ping()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("DB: Successfully connected")
// }

func connectWithRetry(cfg database.PostgresSQLConfig, maxAttemps int, deleay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for attemps := 1; attemps <= maxAttemps; attemps++ {
		log.Printf("PostgresSQL connection attemp %d/%d\n", attemps, maxAttemps)

		db, err = database.NewPostgresSQLStorage(cfg)
		if err == nil {
			log.Println("PostgresSQL conection successfully complete")
			return db, nil
		}

		log.Printf("Attemp %d/%d failed: %v\n", attemps, maxAttemps, err)

		if attemps < maxAttemps {
			log.Printf("Wating %v before the next conection attemp", deleay)
			time.Sleep(deleay)
		}
	}

	return nil, fmt.Errorf("all connection atemps failed: %w", err)
}

func initDatabase(db *sql.DB) error {
	log.Println("Initializing database structure")

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error on verify conection: %w", err)
	}

	log.Println("Successfully verify DB conection")
	return nil
}

func shutdown(server *api.APIServer) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Recived Signal: %v. Initializing shutdown.\n", sig)
	log.Println("Application has ended.")
	os.Exit(0)
}

func getEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		return "desenvolvimento"
	}
	return env
}
