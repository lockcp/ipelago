package main

import "github.com/labstack/echo/v4"

func main() {
	defer db.DB.Close()

	e := echo.New()
	e.HTTPErrorHandler = errorHandler

	e.Static("/public", "public")

	e.File("/", "public/my-island-info.html")

	api := e.Group("/api", sleep)
	api.GET("/get-my-island", getMyIsland)
	api.POST("/create-my-island", createMyIsland)
	api.GET("/my-messages", myMessages)

	e.Logger.Fatal(e.Start(*addr))
}
