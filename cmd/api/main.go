package main

import (
	"MiniJira/internal/config"
	"MiniJira/internal/httpapi"
	"MiniJira/internal/store/memory"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s := memory.NewStore()
	cfg := config.LoadConfig()

	mux := httpapi.NewMux(s)

	srv := &http.Server{Addr: ":" + cfg.HTTPPort, Handler: mux}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err)
			}
			return
		}

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = srv.Shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}

		return
	}()
}
