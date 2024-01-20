package main

import (
	"fmt"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"chatserver/internal/utils"
	"chatserver/internal/server"
)

func main() {
	jwtKeyFilename := utils.GetEnvVar("JWTKEY_FILENAME", "public_key.der")
	serverHost := utils.GetEnvVar("SERVER_HOST", "localhost")
	serverPort := utils.GetEnvVar("SERVER_PORT", "9000")
	dbHost := utils.GetEnvVar("POSTGRES_HOST", "localhost")
	dbPort := utils.GetEnvVar("POSTGRES_PORT", "5432")
	dbName := utils.GetEnvVar("POSTGRES_DB", "gochat")
	dbUser := utils.GetEnvVar("POSTGRES_USERNAME", "gochat")
	dbPassword := utils.GetEnvVar("POSTGRES_PASSWORD", "gochat")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
		return
	}

	chatserver, err := server.NewChatserver(serverHost, serverPort, db, jwtKeyFilename)
	if err != nil {
		log.Fatalln(err)
		return
	}
	chatserver.Run()
}