package main

import "github.com/labstack/echo"

func health(e echo.Context) error  {
	e.String(200, "ok")
	return nil
}
