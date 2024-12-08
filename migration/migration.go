package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/petersizovdev/MEDODS-T.git/pkg/env"
)

func main() {
	direction := flag.String("direction", "up", "Direction of migration: up or down")
	flag.Parse()

	err := env.LoadEnv(".env")
	if err != nil {
		fmt.Println("Err to load .env", err)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	m, err := migrate.New("file://migration/migrations", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	switch *direction {
	case "up":
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				fmt.Println("Нет изменений")
			} else {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Миграция прошла успешно")
		}
	case "down":
		if err := m.Down(); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Откат миграции прошел успешно")
		}
	default:
		log.Fatalf("Unknown direction: %s", *direction)
	}
}
