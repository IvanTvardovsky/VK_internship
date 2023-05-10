package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"vk_telegram_bot/internal/config"
)

func Init(config *config.Config) *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DataBase.Host, config.DataBase.Port, config.DataBase.Username, config.DataBase.Password, config.DataBase.Database)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
