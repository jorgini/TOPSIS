package usecase

import (
	"webApp/app/configs"
	"webApp/app/repository"
)

type DiContainer struct {
	singleton configs.PostgresInstance
	config    *configs.Config
}

func NewDiContainer(config *configs.Config) *DiContainer {
	return &DiContainer{
		singleton: configs.ConnectToDb(&config.PgClientCfg),
		config:    config,
	}
}

func (di *DiContainer) GetInstanceRepository() *repository.Repository {
	return repository.NewRepository(di.singleton.DB, &di.config.DbConfig)
}

func (di *DiContainer) GetInstanceService() *Service {
	return NewService(di.GetInstanceRepository())
}

func (di *DiContainer) ShutDown() error {
	return di.singleton.DB.Close()
}
