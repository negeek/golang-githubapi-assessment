package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	githubCrons "github.com/negeek/golang-githubapi-assessment/crons/v1/github"
	"github.com/negeek/golang-githubapi-assessment/db"
	"github.com/negeek/golang-githubapi-assessment/middlewares"
	routes "github.com/negeek/golang-githubapi-assessment/routes/v1"
)

func loadEnv() {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "dev" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("error loading .env file: ", err)
		}
	}
}

func setupDB() {
	dbUrl := os.Getenv("DATABASE_URL")
	log.Println("connecting to DB...")
	if err := db.Connect(dbUrl); err != nil {
		log.Fatal("failed to connect to DB: ", err)
	}
	log.Println("connected to DB")
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.CORS)
	routes.V1routes(router.StrictSlash(true))
	return router
}

func setupServer() *http.Server {
	router := setupRouter()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("start server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Server error: ", err)
		}
	}()

	return server
}

func setupCronJobs() {
	githubCrons.AddFunc("@hourly", githubCrons.CommitCron)
	githubCrons.Start()
}

func main() {
	loadEnv()
	setupDB()
	setupCronJobs()

	server := setupServer()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown Error: ", err)
	}

	log.Println("Server gracefully stopped")
}
