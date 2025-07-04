package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresSQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresSQLStorage(cfg PostgresSQLConfig) (*sql.DB, error) {
	// Formato: postgres://user:password@host:port/dbname?sslmode=disable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro in opennig a conection to the Postgres DB: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetConnMaxIdleTime(5)
	// db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error on conecting to PostgresSQL: %w", err)
	}

	log.Println("Connection with PostgresSQL successfully complete. ")
	return db, nil
}

// CreateTables cria as tabelas necessárias para o e-commerce
func CreateTables(db *sql.DB) error {
	// Exemplo de criação de tabela de usuários
	// PostgreSQL tem tipos de dados específicos como SERIAL para IDs auto-incrementais
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		// Adicione mais tabelas conforme necessário
		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2) NOT NULL,
			stock INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	// Executando cada query de criação
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("erro ao criar tabela: %w", err)
		}
	}

	log.Println("Tabelas criadas com sucesso!")
	return nil
}
