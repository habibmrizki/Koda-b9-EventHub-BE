package main

import (
	"log"

	_ "github.com/habibmrizki/BE-EventHub/docs"
	"github.com/habibmrizki/BE-EventHub/internal/configs"
	"github.com/habibmrizki/BE-EventHub/internal/routes"

	"github.com/joho/godotenv"
)

// @title 											Event-Hub
// @version 										1.0
// @description 									Event-Hub
// @host											localhost:3000
// @securityDefinitions.apikey 	JWTtoken
// @in header
// @name Authorization
func main() {
	// Inisialisasi Database for this project
	if err := godotenv.Load(); err != nil {
		log.Println(err.Error())
		return
	}

	db, err := configs.InitDB()
	if err != nil {
		log.Println("FAILED TO CONNECT DB")
		return
	}

	defer db.Close()

	err = configs.PingDB(db)
	if err != nil {
		log.Println("PING TO DB FAILED", err.Error())
		return
	}

	log.Println("DB CONNECTED")

	// JALANKAN SERVER GIN
	r := routes.InitRouter(db)

	r.Run(":3000")
}
