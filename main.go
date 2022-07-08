package main

import (
	"time"
	"website/internal/app"

	"github.com/labstack/echo/v4"
)

func main() {
	ch := make(chan *echo.Echo, 1)

	go func() {
		e := app.Run()
		ch <- e
		e.Start(":8080")
	}()

	for {
		time.Sleep(10 * time.Minute)
		go func() {
			c := <-ch
			e := app.Run()
			ch <- e
			c.Close()
			e.Start(":8080")
		}()
	}
}
