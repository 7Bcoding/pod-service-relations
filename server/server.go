package server

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"pod-service-relations/config"
	"pod-service-relations/logging"
)

var g errgroup.Group

/*
Init
1. start web server
*/
func Init() {
	conf := config.NewServerConfig()
	logger := logging.GetLogger()

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", conf.Bind, conf.Port),
		Handler:      newRoute(conf),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	// TODO start manager

	// start web server
	g.Go(func() error {
		logger.Infof("http: starting web server at %s", server.Addr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		logger.Fatal(err)
	}
}
