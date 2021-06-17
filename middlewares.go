package main

import (
	"time"

	"github.com/labstack/echo/v4"
)

func sleep(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		time.Sleep(time.Second)
		return next(c)
	}
}
