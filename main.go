package main

import (
	"context"
	"log"
	"time"
	"website/internal/app"
	"website/pkg/config"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

var (
	l   *log.Logger
	ctx context.Context
)

func main() {
	l = log.Default()
	ctx = context.Background()
	if err := config.Init("config"); err != nil {
		l.Fatalf("init config failed %s", err)
		return
	}

	ch := make(chan *echo.Echo, 1)

	viper.SetDefault("server.update.duration", "5m")
	duration := viper.GetDuration("server.update.duration")

	go func() {
		e := app.Run()
		ch <- e
		// e.Start(":8080")
		e.StartAutoTLS(":8080")
	}()

	for {
		time.Sleep(duration)
		go func() {
			c := <-ch
			e := app.Run()
			ch <- e
			c.Close()
			// e.Start(":8080")
			e.StartAutoTLS(":8080")
		}()
	}
}
