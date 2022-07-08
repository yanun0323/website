package app

import (
	"website/internal/repository"
	"website/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() *echo.Echo {

	e := echo.New()
	e.Logger.SetLevel(4)

	svc := service.NewService(repository.NewRepo())

	rateLimiter := middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))
	svc.SetHomePage(e, rateLimiter)
	svc.SetAllArticlePage(e, rateLimiter)

	return e
}
