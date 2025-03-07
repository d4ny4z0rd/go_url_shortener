package main

import (
	"log"
	"url_shortener_1/db"
	"url_shortener_1/env"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading env variables")
	}

	db, err := db.NewDB(db.DBConfig{
		User:     env.GetString("DB_USER", "admin"),
		Password: env.GetString("DB_PASSWORD", "password"),
		Host:     env.GetString("DB_HOST", "localhost"),
		Port:     env.GetString("DB_PORT", "5432"),
		DBName:   env.GetString("DB_NAME", "project"),
		SSLmode:  env.GetString("SSL_MODE", "disable"),
	})
	if err != nil {
		log.Fatal(err)
	}

	app := NewApplication(env.GetString("ADDR", ":8080"), db)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
	