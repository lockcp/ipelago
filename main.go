package main

import "github.com/labstack/echo/v4"

func main() {
	defer db.DB.Close()

	e := echo.New()
	e.HTTPErrorHandler = errorHandler

	e.Static("/public", "public")

	api := e.Group("/api")
	api.GET("/get-my-island", getMyIsland)
	api.POST("/create-my-island", createMyIsland)

	e.Logger.Fatal(e.Start(*addr))
}
