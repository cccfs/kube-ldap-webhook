package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func other(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"content": "ok",
	})
}
