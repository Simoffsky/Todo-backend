package todoapp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"todo/internal/config"
	"todo/internal/repository"
	"todo/internal/repository/task"
	"todo/pkg/log"

	"github.com/go-redis/redis"
)

type App struct {
	config *config.Config
	logger log.Logger

	authService AuthService
	taskService TaskService
}

func NewApp(config *config.Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Run() error {
	err := a.configureApp()
	if err != nil {
		return err
	}

	return a.startHTTPServer()
}

func (a *App) configureApp() error {
	a.logger = log.NewDefaultLogger(
		log.LevelFromString(a.config.LoggerLevel),
	).WithTimePrefix(time.DateTime)

	a.logger.Debug("Config parameters: " + fmt.Sprintf("%+v", a.config))

	var err error
	a.authService, err = NewAuthServiceGRPC(a.config.AuthAddr, a.config.JwtSecret)
	if err != nil {
		return err
	}

	db, err := repository.ConnectToDB(a.config.DbConn)
	if err != nil {
		return err
	}
	a.logger.Info("Connected to DB")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     a.config.RedisAddr,
		Password: "",
		DB:       0,
	})

	err = redisClient.Ping().Err()

	if err != nil {
		return err
	}

	a.logger.Info("Connected to Redis")

	taskRepo := task.NewPostgresTaskRepository(db)
	taskListRepo := task.NewPostgresTaskListRepository(db)
	a.taskService = NewTaskService(taskRepo, taskListRepo, redisClient)
	return nil
}

func (a *App) startHTTPServer() error {
	mux := http.NewServeMux()

	a.registerHandlers(mux)

	server := &http.Server{
		Addr:         ":" + a.config.ServerPort,
		Handler:      mux,
		ReadTimeout:  a.config.HTTPTimeout,
		WriteTimeout: a.config.HTTPTimeout,
	}

	return a.gracefulStart(server)
}

func (a *App) registerHandlers(mux *http.ServeMux) {
	mux.Handle("/task", a.ProtectMiddleware(http.HandlerFunc(a.handleCreateTask)))
	mux.Handle("/task/{task_id}", a.ProtectMiddleware(http.HandlerFunc(a.handleTask)))

	mux.Handle("/task-list", a.ProtectMiddleware(http.HandlerFunc(a.handleCreateTaskList)))
	mux.Handle("/task-list/{task_id}", a.ProtectMiddleware(http.HandlerFunc(a.handleTaskList)))
	mux.HandleFunc("/register", a.handleRegister)
	mux.HandleFunc("/login", a.handleLogin)
}

// Starts HTTP Server with graceful shutdown
func (a *App) gracefulStart(server *http.Server) error {
	errCh := make(chan error, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		a.logger.Info("Server started on port " + a.config.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-quit:
		a.logger.Info("Server is shutting down...")
	case err := <-errCh:
		a.logger.Error("HTTP Server error:" + err.Error())
	}

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return server.Shutdown(ctx)
}
