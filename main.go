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

	githubFuncs "github.com/negeek/golang-githubapi-assessment/api/v1/github"
	githubCrons "github.com/negeek/golang-githubapi-assessment/crons/v1/github"
	githubModels "github.com/negeek/golang-githubapi-assessment/data/v1/github"
	"github.com/negeek/golang-githubapi-assessment/db"
	"github.com/negeek/golang-githubapi-assessment/middlewares"
	routes "github.com/negeek/golang-githubapi-assessment/routes/v1"
	runners "github.com/negeek/golang-githubapi-assessment/runners/v1/github"
)

func loadEnv() {
	log.Println(("load env"))
	environment := os.Getenv("ENVIRONMENT")
	if environment == "dev" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("error loading .env file: ", err)
		}
	}
}

func connectDB() {
	dbUrl := os.Getenv("POSTGRESQL_URL")
	log.Println("connecting to DB...")
	if err := db.Connect(dbUrl); err != nil {
		log.Fatal("failed to connect to DB: ", err)
	}
	log.Println("connected to DB")
}

func seedDB() {
	log.Println("seed db with data")
	githubModels.SetEnvSetupData()
}

func setupRouter() *mux.Router {
	log.Println("setup router")
	router := mux.NewRouter()
	router.Use(middlewares.CORS)
	routes.V1routes(router.StrictSlash(true))
	return router
}

func startServer() *http.Server {
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
			log.Println("server error: ", err)
		}
	}()

	return server
}

func startServices() {
	log.Println("start services")
	go func() {
		log.Println("start github api data fetching service")
		runners.Run(1, githubFuncs.GithubHandler)
	}()
}

func startCronJobs() {
	log.Println("setup cronjobs")
	githubCrons.AddJob("@hourly", githubFuncs.GithubHandler)
	githubCrons.StartCron()
}

func main() {
	loadEnv()
	connectDB()
	seedDB()
	startServices()
	startCronJobs()

	server := startServer()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("server shutdown error: ", err)
	}

	log.Println("server gracefully stopped")
}
