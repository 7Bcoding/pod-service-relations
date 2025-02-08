// Copyright 2022 Baidu Inc. All rights reserved.
// Use of this source code is governed by a xxx
// license that can be found in the LICENSE file.

// Package main is special.  It defines a
// standalone executable program, not a library.
// Within package main the function main is also
// special—it’s where execution of the program begins.
// Whatever main does is what the program does.
package main

import (
	"context"
	"pod-service-relations/config"
	"pod-service-relations/database"
	"pod-service-relations/logging"
	"pod-service-relations/server"
	"pod-service-relations/service"
)

// main the function where execution of the program begins
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config.InitConfigs()
	logging.Init()
	database.Init()
	defer database.Close()
	service.GetPodToServiceRelations(ctx)
	server.Init()
}
