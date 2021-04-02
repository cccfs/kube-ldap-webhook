package main

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"os"
)

func main()  {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	if _, err := os.Stat(".env"); !errors.Is(err, os.ErrNotExist) {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Load .env files error")
		}
	}
	authorized := e.Group("/auth", middleware.BasicAuth(auth))
	{
		authorized.GET("/other", other)
	}
	e.GET("/health", health)
	e.POST("/token", token)
	e.Logger.Fatal(e.Start(":8080"))
}
