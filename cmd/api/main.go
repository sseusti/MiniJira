// @Title MiniJira API
// @Version 0.1
// @Description ...
// @BasePath /
package main

import (
	"MiniJira/internal/config"
	"MiniJira/internal/httpapi"
	"MiniJira/internal/store/memory"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	s := memory.NewStore()
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.New()
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}
	logger.SetLevel(level)
	if cfg.LogFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	logger.Info("starting server")

	mux := httpapi.NewMux(s, s, s, logger)

	srv := &http.Server{Addr: ":" + cfg.HTTPPort, Handler: mux}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			logger.WithError(err).Fatal("error starting server")
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		logger.WithError(err).Fatal("error shutting down server")
	}

	return
}
