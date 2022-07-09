package main

import (
	"log"
	"time"
	"website/internal/app"
	"website/pkg/config"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

var (
	LOG *log.Logger
)

func main() {
	LOG = log.Default()
	if err := config.Init("config"); err != nil {
		LOG.Fatalf("init config failed %s", err)
		return
	}
	crt := viper.GetString("server.crt")
	key := viper.GetString("server.key")

	ch := make(chan *echo.Echo, 1)

	go func() {
		e := app.Run()
		ch <- e
		LOG.Fatal(e.StartTLS(":443", []byte(crt), []byte(key)))
	}()

	for {
		time.Sleep(10 * time.Minute)
		go func() {
			c := <-ch
			e := app.Run()
			ch <- e
			LOG.Fatal(c.Close())
			LOG.Fatal(e.StartTLS(":443", []byte(crt), []byte(key)))
		}()
	}
}
