package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fuadop/altschool-capstone-scissor/model"
	"github.com/labstack/echo/v4"
)

type SBody struct {
	URL string `json:"url"` 
}
type SRes struct {
	ID string `json:"id"`
	// shortened URL
	URL string `json:"url"`
}

//	@Summary		Healty status checker
//	@Description	Healty status checker for load balancers and monitoring systems.
//	@Tags			api
//	@Produce		json
//	@Success		200	{object}	JSONResponse[any]
//	@Failure		400	{object}	JSONResponse[any]
//	@Failure		500	{object}	JSONResponse[any]
//	@Router			/api/health [get]
func Health(c echo.Context) error {
	return HandleResponseJSON(c, 200, "OK", nil)
}

//	@Summary		Shorten a URL
//	@Description	Shorten a URL
//	@Tags			api
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		SBody	true	"Request body"
//	@Success		200		{object}	JSONResponse[SRes]
//	@Failure		400		{object}	JSONResponse[any]
//	@Failure		500		{object}	JSONResponse[any]
//	@Router			/api/shorten [post]
func Shorten(c echo.Context) error {
	var body SBody
	if err := c.Bind(&body); err != nil {
		return HandleResponseJSON(c, 422, err.Error(), nil)
	}
	url, err := url.Parse(body.URL)
	if err != nil {
		return HandleResponseJSON(c, 422, err.Error(), nil)
	}	
	if (url.Scheme != "http" && url.Scheme != "https") || url.Hostname() == "" {
		return HandleResponseJSON(c, 422, "Validation Error: URL is invalid, required format: http://berlin.de", nil)
	}

	res, err := http.Get(body.URL)
	if err != nil {
		msg := fmt.Sprintf("URL not responding. %s", err.Error())
		return HandleResponseJSON(c, 422, msg, nil)
	}
	if res.StatusCode < 200 || res.StatusCode >= 400 {
		msg := fmt.Sprintf("URL not OK. respnded with status %d", res.StatusCode)
		return HandleResponseJSON(c, 422, msg, nil)
	}

	id, err := model.URLIndex(body.URL)
	if err != nil {
		return HandleResponseJSON(c, http.StatusInternalServerError, err.Error(), nil)
	}

	// * resolve domain name
	domainName, port := os.Getenv("DOMAIN_NAME"), os.Getenv("PORT")
	if domainName == "" || domainName == "localhost" {
		if port == "" {
			port = "8080"
		}
		domainName = fmt.Sprintf("http://localhost:%s", port)
	}
	if strings.HasPrefix(domainName, "http") != true {
		domainName = fmt.Sprintf("https://%s", domainName)
	}

	data := map[string]string{
		"id": id,
		"url": fmt.Sprintf("%s/%s", domainName, id),
	}
	return HandleResponseJSON(c, http.StatusCreated, "URL shortened", data)
}

//	@Summary		Unpublish/Delete a shortened URL
//	@Description	Unpublish/Delete a shortened URL.
//	@Tags			api
//	@Produce		json
//	@Param			id	path		string	true	"The URL ID"
//	@Success		200	{object}	JSONResponse[any]
//	@Failure		400	{object}	JSONResponse[any]
//	@Failure		404	{object}	JSONResponse[any]
//	@Failure		500	{object}	JSONResponse[any]
//	@Router			/api/unpublish/{id} [delete]
func Unpublish(c echo.Context) error {
	id := c.Param("id")
	if err := model.UnpublishIndex(id); err != nil {
		return HandleResponseJSON(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return HandleResponseJSON(c, http.StatusOK, "URL unpublished", nil)
}

