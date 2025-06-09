package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Bayan2019/hackathon-2025-police-volunteers/configuration"
	"github.com/joho/godotenv"
)

// @title VOLUNTEERS API
// @version 1.0
// @description This is a sample server POLICE.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host hackathon-2025-zfmi.onrender.com
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("warning: assuming default configuration. .env unreadable: %v\n", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	// fmt.Println(dbURL)
	err = configuration.Connect2DB(dbURL)
	if err != nil {
		log.Println("DATABASE_URL environment variable is not set")
		log.Println("Running without CRUD endpoints")
		fmt.Println(err.Error())
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "hackathon2025"
	}

	if configuration.ApiCfg != nil {
		configuration.ApiCfg.JwtSecret = jwtSecret
	} else {
		fmt.Println("No DATABASE_URL")
		configuration.ApiCfg = &configuration.ApiConfiguration{
			JwtSecret: jwtSecret,
		}
	}

}
