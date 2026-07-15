package core

import (
	"errors"
	"fmt"

	"github.com/0xMoonrise/gochive/internal/config"
)

type Stage func(*App) error

var (
	StageStorage Stage = (*App).InitStorage
	StageDB      Stage = (*App).InitDB
	StageConfig  Stage = (*App).InitConfig
)

func NewApp() *App {
	return &App{}
}

func (a *App) Run(stages ...Stage) error {
	for _, stage := range stages {
		if err := stage(a); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) InitConfig() error {
	conf, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	a.Config = conf
	return nil
}

func (a *App) InitStorage() error {
	if a.Config == nil {
		return errors.New("a config must be set before init Storage")
	}
	client, err := a.SetMode()
	if err != nil {
		return fmt.Errorf("creating storage client: %w", err)
	}
	a.Storage = client
	return nil
}

func (a *App) InitDB() error {
	if a.Config == nil {
		return errors.New("a config must be set before init DB")
	}
	closeDB, err := BootDatabase(a)
	if err != nil {
		return fmt.Errorf("booting database: %w", err)
	}
	a.closeDB = closeDB
	return nil
}

func (a *App) Cleanup() error {
	if a.closeDB == nil {
		return nil
	}
	return a.closeDB()
}
