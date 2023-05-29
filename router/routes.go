package router

import (
	"fmt"
	"os"

	"github.com/fuadop/altschool-capstone-scissor/controller"
	"github.com/labstack/echo/v4"

	"github.com/fuadop/altschool-capstone-scissor/docs"
	"github.com/swaggo/echo-swagger"
)

//	@title			Scissors
//	@version		1.0
//	@description	Fast minimalist URL shortener

//	@contact.name	Fuad Olatunji
//	@contact.url	https://fuadolatunji.me
//	@contact.email	fuadolatunji@gmail.com

//	@Schemes	http https
//	@BasePath	/
func Register(e *echo.Echo) {
	// * set dynamic swagger informations 
	domainName, port := os.Getenv("DOMAIN_NAME"), os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if domainName == "" {
		domainName = "localhost"
	}

	docs.SwaggerInfo.Host = domainName
	if domainName == "localhost" {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", port)
	}

	// * register app routes
	e.GET("/:id", controller.Redirect)

	api := e.Group("/api") // * /api	
	api.GET("/health", controller.Health)
	api.POST("/shorten", controller.Shorten)
	api.DELETE("/unpublish/:id", controller.Unpublish)

	api.GET("/docs/*", echoSwagger.WrapHandler)
	api.GET("/swagger/*", echoSwagger.WrapHandler)

	analytics := api.Group("/analytics") // * /api/analytics
	analytics.GET("/:id", controller.URLAnalytics)
}

