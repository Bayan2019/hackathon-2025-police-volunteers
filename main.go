package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Bayan2019/hackathon-2025-police-volunteers/configuration"
	"github.com/Bayan2019/hackathon-2025-police-volunteers/controllers"
	_ "github.com/Bayan2019/hackathon-2025-police-volunteers/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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

// @host hackathon-2025-police-volunteers.onrender.com
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

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/", controllers.HelloHandler)

	router.Get("/swagger/*",
		httpSwagger.Handler(httpSwagger.URL("https://hackathon-2025-police-volunteers.onrender.com/swagger/doc.json")))

	v1Router := chi.NewRouter()

	if configuration.ApiCfg.DB != nil {
		authHandlers := controllers.NewAuthHandlers(configuration.ApiCfg.DB, configuration.ApiCfg.JwtSecret)

		v1Router.Post("/auth/sign-in", authHandlers.Login)
		v1Router.Post("/auth/refresh", authHandlers.Refresh)
		v1Router.Post("/auth/sign-out", authHandlers.Logout)
	}

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: time.Second * 10,
	}

	log.Printf("Serving on: http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())

}
