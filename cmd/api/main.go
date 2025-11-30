package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/arunhiremath92/go-shorturl/internal/api"
	"github.com/arunhiremath92/go-shorturl/internal/urlshortner"
	"github.com/arunhiremath92/go-shorturl/pkg/redis"
)

func main() {
	loghandler := log.New(os.Stdout, "url-shortner:", log.Default().Flags())
	addr := ":8000"

	redisAddr := os.Getenv("REDIS_ADDRESS")
	if redisAddr == "" {
		redisAddr = "redis:6379" // For local development
	}
	redisConfig := redis.RedisStoreConfig{
		Address:   redisAddr,
		DefaultDb: 0,
	}

	redisStore := redis.NewRedisStore(redisConfig)
	urlShortenerSvc := urlshortner.NewUrlShortner(redisStore)

	app := api.NewApp(api.AppConfig{Addr: addr, Logger: loghandler,
		UrlService: urlShortenerSvc})

	// handle errors from the http
	appErrorHandler := make(chan error, 1)
	// setup interruption Handler
	signalHandler := make(chan os.Signal, 1)
	signal.Notify(signalHandler, os.Interrupt)

	go func() {
		appErrorHandler <- app.Start()
	}()

	select {
	case serverError := <-appErrorHandler:
		loghandler.Printf("server error: %v", serverError)
	case <-signalHandler:
		loghandler.Print("received interruption from the os")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := app.ShutDownServer(ctx); err != nil {
			loghandler.Printf("graceful shutdown failed: %v", err)
		}
	}

}
