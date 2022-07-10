package app

import (
	"fmt"
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
	m := []echo.MiddlewareFunc{rateLimiter, middleware.WWWRedirect(), middleware.HTTPSRedirect()}
	fmt.Println(m)
	svc.SetHomePage(e, m...)
	svc.SetAllArticlePage(e, m...)

	return e
}
