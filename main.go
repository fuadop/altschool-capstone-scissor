package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fuadop/altschool-capstone-scissor/router"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	godotenv.Load()
}

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	router.Register(e)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	log.Fatalln(e.Start(addr))
}
