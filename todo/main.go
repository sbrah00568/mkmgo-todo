package main

import (
	"context"
	"fmt"
	"mkmgo-todo/todo/handler"
	"mkmgo-todo/todo/task"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Setup database
	db, err := gorm.Open(sqlite.Open("todo/gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Database connection failed")
	}
	db.AutoMigrate(&task.Task{})

	// Setup repository, service, and handlers
	taskRepo := task.NewTaskRepositoryImpl(db)
	taskSvc := task.NewTaskServiceImpl(taskRepo)
	taskHandler := handler.NewTaskHandler(taskSvc)

	handler := Handler{taskHandler: taskHandler}

	// Setup router and server
	router := mux.NewRouter()
	setupRoutes(router, handler)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	log.Info().Msg("Start server")
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
		fmt.Printf("Starting server at port 8080")
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server shutdown failed")
	}

	log.Info().Msg("Server stopped successfully")
}

type Handler struct {
	taskHandler *handler.TaskHandler
}

func setupRoutes(router *mux.Router, h Handler) {
	router.HandleFunc("/todo/tasks/health", healthCheck).Methods("GET")
	router.HandleFunc("/todo/tasks", h.taskHandler.WriteTaskHandler).Methods("POST")
	router.HandleFunc("/todo/tasks/{id}", h.taskHandler.UpdateTaskHandler).Methods("PATCH")
	router.HandleFunc("/todo/tasks", h.taskHandler.GetAllTaskHandler).Methods("GET")
	router.HandleFunc("/todo/tasks/{id}", h.taskHandler.DeleteTaskHandler).Methods("DELETE")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
