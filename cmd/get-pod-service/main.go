package main

import (
	"context"
	"pod-service-relations/config"
	"pod-service-relations/database"
	"pod-service-relations/logging"
	"pod-service-relations/service"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config.InitConfigs()
	logging.Init()
	database.Init()
	service.GetPodToServiceRelations(ctx)
}
