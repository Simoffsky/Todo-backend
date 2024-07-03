package todoapp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
	"todo/internal/config"
	"todo/pkg/log"
)

type App struct {
	config *config.Config
	logger log.Logger
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

	return nil
}

func (a *App) startHTTPServer() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			a.logger.Error("Error writing response: " + err.Error())
		}
	})

	server := &http.Server{
		Addr:         ":" + a.config.ServerPort,
		Handler:      mux,
		ReadTimeout:  a.config.HTTPTimeout,
		WriteTimeout: a.config.HTTPTimeout,
	}

	return a.gracefulStart(server)
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
