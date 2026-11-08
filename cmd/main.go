package main

import (
	"arshin/internal/adapter/sqlite"
	"arshin/internal/adapter/xlsx"
	"arshin/internal/adapter/xml"
	"arshin/internal/app"
	"arshin/internal/gui"
	"database/sql"
	"io"
	"log"

	_ "github.com/gwenn/gosqlite"
)

func main() {
	db, err := sql.Open("sqlite3", "/home/mitsuha/test.db")
	if err != nil {
		log.Fatal(err)
	}

	pr, pw := io.Pipe()
	configRepo := sqlite.New(db)
	xmlConverter := xml.NewConverter()
	logger := log.New(pw, "ОШИБКА КОНВЕРТАЦИИ:", 0)
	xlsxParser := xlsx.NewParser(logger)
	service := app.New(configRepo, logger)

	gui.NewApp(xlsxParser, xmlConverter, xmlConverter, "testfgisurl", pr, service).Run()
}
