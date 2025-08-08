package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	HOST     string
	PORT     int
	USER     string
	PASSWORD string
	DB_NAME  string
}

type Database struct {
	DB *sql.DB
}

func BuildDataBaseConfig() (*DatabaseConfig, error) {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, continuing...")
	}

	viper.SetConfigFile("config/config.json")
	viper.SetConfigType("json")

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &DatabaseConfig{
		HOST:     viper.GetString("DATABASE.HOST"),
		PORT:     viper.GetInt("DATABASE.PORT"),
		USER:     viper.GetString("DATABASE.USER"),
		DB_NAME:  viper.GetString("DATABASE.DB_NAME"),
		PASSWORD: os.Getenv("DATABASE_PASSWORD"),
	}, nil
}

func InitDB(dbConfig *DatabaseConfig) (*Database, error) {
	dsn := dbConfig.GetDatabaseDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")

	return &Database{DB: db}, nil
}

func (c *DatabaseConfig) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", c.HOST, c.PORT, c.USER, c.PASSWORD, c.DB_NAME)
}

func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

func (d *Database) Ping() error {
	return d.DB.Ping()
}

func (d *Database) Begin() (*sql.Tx, error) {
	return d.DB.Begin()
}

func (d *Database) Prepare(query string) (*sql.Stmt, error) {
	return d.DB.Prepare(query)
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.DB.Query(query, args...)
}

func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRow(query, args...)
}
