package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"os"
	"time"
)

type Server struct {
	WebServer *http.Server
	Router    *chi.Mux
	Gorm      *gorm.DB
	DB        *sql.DB
}

func (s *Server) InitDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Get generic database object
	d, err := db.DB()
	if err != nil {
		log.Errorln(err)
	}

	// Validate DB DSN data:
	err = d.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	// Assign database connection to Server
	s.Gorm = db
	s.DB = d
}

func (s *Server) InitWebServer() {
	// Init router
	s.Router = chi.NewRouter()

	// Apply middlewares to router
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(60 * time.Second))

	// Prepare webserver
	s.WebServer = &http.Server{
		Handler:      s.Router,
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("HTTP_URL"), os.Getenv("HTTP_PORT")),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	// Init the routes
	InitRoutes()
}

func (s *Server) Init() {

	// Initialize database
	s.InitDatabase()

	// Initialize database
	s.InitWebServer()
}

func (s *Server) Run() {
	log.Infof("http server running on %s:%s\n", os.Getenv("HTTP_URL"), os.Getenv("HTTP_PORT"))
	if err := srv.WebServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}

func (s *Server) Stop() {
	// Close database connection
	if err := srv.DB.Close(); err != nil {
		log.Fatalln(err)
	}
	log.Infoln("Database connection closed")

	// Stop http server
	srv.WebServer.SetKeepAlivesEnabled(false)
	if err := srv.WebServer.Shutdown(context.Background()); err != nil {
		log.Fatalln(err)
	}
	log.Infoln("Http server stopped")
}
