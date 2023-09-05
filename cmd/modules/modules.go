package modules

import (
	metrics "github.com/meowt/SberTestTask/internal/metrics/handler"
	objects "github.com/meowt/SberTestTask/internal/objects/handler"
	objectsRepo "github.com/meowt/SberTestTask/internal/objects/repo"
	probes "github.com/meowt/SberTestTask/internal/probes/handler"
	"github.com/meowt/SberTestTask/internal/storage"
)

func Setup(s *storage.Storage) *HandlerModule {
	return SetupHandlerModule(SetupRepoModule(s))
}

type RepoModule struct {
	objects.RepoService
}

func SetupRepoModule(s *storage.Storage) RepoModule {
	return RepoModule{
		RepoService: objectsRepo.SetupObjectsRepo(s),
	}
}

type HandlerModule struct {
	MetricsHandler metrics.Handler
	ObjectsHandler objects.Handler
	ProbesHandler  probes.Handler
}

func SetupHandlerModule(repo RepoModule) *HandlerModule {
	return &HandlerModule{
		MetricsHandler: metrics.SetupMetricsHandler(),
		ObjectsHandler: objects.SetupObjectsHandler(repo),
		ProbesHandler:  probes.SetupProbesHandler(),
	}
}
