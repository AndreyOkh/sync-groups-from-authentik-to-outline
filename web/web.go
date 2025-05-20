package web

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"sync-groups-from-authentik-to-outline/authentik"
	"sync-groups-from-authentik-to-outline/config"
	"sync-groups-from-authentik-to-outline/outline"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

var isShuttingDown atomic.Bool

type WebserverDeps struct {
	OClient *outline.Client
	AClient *authentik.Client
	Config  *config.Conf
}

func Webserver(deps *WebserverDeps) {
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := http.NewServeMux()
	NewHandler(router, Handler{
		OClient:       deps.OClient,
		AClient:       deps.AClient,
		GroupPrefix:   deps.Config.App.GroupPrefix,
		GroupSelector: deps.Config.App.GroupNameSelector,
	})

	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	server := http.Server{
		Addr: ":8081",
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
		Handler: router,
	}
	go func() {
		log.Println("Server starting on :8081")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-rootCtx.Done()
	stop()
	isShuttingDown.Store(true)
	log.Println("Received shutdown signal, shutting down.")

	time.Sleep(_readinessDrainDelay)
	log.Println("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()
	err := server.Shutdown(shutdownCtx)
	stopOngoingGracefully()
	if err != nil {
		log.Println("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(_shutdownHardPeriod)
	}

	log.Println("Server shut down gracefully.")
}
