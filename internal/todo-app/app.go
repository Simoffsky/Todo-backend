package todoapp

import (
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

	a.logger.Debug("Config parameters: " + a.config.String())

	return nil
}

func (a *App) startHTTPServer() error {
	return nil
}
