package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello, world")
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", hello)
	e.File("/favicon.ico", "static/images/favicon.ico")
	e.Logger.Fatal(e.Start(":1323"))
}
