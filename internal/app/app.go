package app

import (
	"website/internal/repository"
	"website/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {

	e := echo.New()
	svc := service.NewService(repository.NewRepo())

	rateLimiter := middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))
	svc.SetHomePage(e, rateLimiter)
	svc.SetAllArticlePage(e, rateLimiter)

	e.Logger.Fatal(e.Start(":8080"))
}
