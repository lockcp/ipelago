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
	api.GET("/get-island/:id", getIslandWithoutMsg)
	api.POST("/create-my-island", createMyIsland)
	api.POST("/update-my-island", updateMyIsland)
	api.GET("/more-my-messages", moreMyMessages)
	api.GET("/more-island-messages", moreIslandMessages)
	api.GET("/more-messages", moreMessages)
	api.POST("/post-message", postMessage)
	api.GET("/publish-newsletter", publishNewsletter)
	api.POST("/follow-island", followIsland)
	api.GET("/all-islands", allIslands)
	api.POST("/update-note", updateNote)
	api.POST("/unfollow", unfollow)
	api.POST("/follow-again", followAgain)
	api.POST("/update-island", updateIsland)

	e.Logger.Fatal(e.Start(*addr))
}
