package main

import (
	"log"
	"time"
	"website/internal/app"
	"website/pkg/config"

	"github.com/labstack/echo/v4"
)

func main() {
	if err := config.Init("config"); err != nil {
		log.Default().Fatal("init config failed")
		return
	}
	// crt := viper.GetString("server.crt")
	// key := viper.GetString("server.key")

	ch := make(chan *echo.Echo, 1)

	go func() {
		e := app.Run()
		ch <- e
		e.StartAutoTLS(":80")
	}()

	for {
		time.Sleep(10 * time.Minute)
		go func() {
			c := <-ch
			e := app.Run()
			ch <- e
			c.Close()
			e.StartAutoTLS(":80")
		}()
	}
}
