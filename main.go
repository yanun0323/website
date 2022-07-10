package main

import (
	"context"
	"log"
	"time"
	"website/internal/app"
	"website/pkg/config"

	"github.com/labstack/echo/v4"
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

	go func() {
		e := app.Run()
		ch <- e
		l.Fatal(e.Start(":8080"))
		// l.Fatal(e.StartAutoTLS(":443"))
	}()

	for {
		time.Sleep(1 * time.Minute)
		go func() {
			c := <-ch
			e := app.Run()
			ch <- e
			l.Fatal(c.Close())
			l.Fatal(e.Start(":8080"))
			// l.Fatal(e.StartAutoTLS(":443"))
		}()
	}
}
