package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

var srv Server
var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetReportCaller(true)
}

func main() {
	// Load dot env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init logging
	initLogger()

	// Initialize the server
	srv.Init()

	// Prepare graceful shutdown
	closed := make(chan struct{})
	go shutdownServer(closed)

	// Run server
	srv.Run()

	<-closed
	log.Println("server stopped")
	os.Exit(0)
}

func shutdownServer(killed chan<- struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	log.Infoln("shutting down ...")

	// Shutdown the server
	srv.Stop()

	close(killed)
}

func initLogger() {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	default:
		log.SetLevel(logrus.ErrorLevel)
	}

	switch os.Getenv("LOG_OUTPUT") {
	case "file":
		// Create log file
		file, err := os.OpenFile("payment.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Out = file
		} else {
			log.Info("Failed to log to file, using default stderr")
		}
	case "stderr":
		log.Out = os.Stderr
	default:
		log.Out = os.Stdout
	}
}
