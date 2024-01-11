package handler

import (
	"github.com/labstack/echo/v4"
	"log"
)

func RequestLogging(fun func(c echo.Context) any) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("记录log")
		return nil
	}
}
