package controller

import (
	"log"
	"net/http"

	"github.com/fuadop/altschool-capstone-scissor/model"
	"github.com/labstack/echo/v4"
)

//	@Summary		Short URL analytics
//	@Description	Fetch analytics of a shortened URL
//	@Tags			api
//	@Produce		json
//	@Param			id	path		string	true	"The URL ID"
//	@Success		200			{object}	JSONResponse[model.URL]
//	@Failure		400			{object}	JSONResponse[any]
//	@Failure		500			{object}	JSONResponse[any]
//	@Router			/api/analytics/{id} [get]
func URLAnalytics(c echo.Context) error {
	id := c.Param("id")
	info, err := model.GetIndex(id)
	if err != nil {
		log.Println(err)
		return HandleResponseJSON(c, http.StatusNotFound, "Index not found", nil)
	}

	return HandleResponseJSON(c, http.StatusOK, "analytics fetched", info)
}

