package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/thenopholo/ecom_backend/cmd/api"
	"github.com/thenopholo/ecom_backend/internal/config"
	"github.com/thenopholo/ecom_backend/internal/database"
	// "golang.org/x/tools/go/analysis/passes/errorsas"
)

func main() {
	db, err := database.NewSQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddr,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

  if err != nil {
    log.Fatal(err)
  }

  initStorage(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}


func initStorage(db *sql.DB) {
  err := db.Ping()
  if err != nil {
    log.Fatal(err)
  }

  log.Println("DB: Successfully connected")
}