package main

import (
	"embed"
	"io"
	"log"
	"megen/models"
	"megen/web"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

//go:embed database.db
var dbFs embed.FS

func main() {
	// create temp file for db
	f, err := os.CreateTemp("", ".db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up

	dbf, err := dbFs.Open("database.db")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = io.Copy(f, dbf)
	if err != nil {
		log.Println(err)
		return
	}
	// opens the database
	db, err := gorm.Open(sqlite.Open(f.Name()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.Debug().AutoMigrate(
		&models.Input{},
		&models.Section{},
	)
	if err != nil {
		log.Println(err)
		return
	}

	web.StartServer(db)
}
