package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fuadop/altschool-capstone-scissor/model"
	"github.com/fuadop/altschool-capstone-scissor/queue"
	"github.com/labstack/echo/v4"
	"github.com/teris-io/shortid"
)

type JSONResponse[K interface{}] struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data K`json:"data"`
}

var jobQueue *queue.JobQueue
func init() {
	jobQueue = queue.NewQueue("default")
}

func Redirect(c echo.Context) error {
	id := c.Param("id")
	url, _ := model.URLFromIndex(id)
	if url == "" {
		return c.String(http.StatusNotFound, "Not Found")
	}

	// schedule job in queue to record analytics
	jobId, err := shortid.Generate()
	if err != nil {
		jobId = "nil"
	}

	ip := c.RealIP()
	log.Println(ip, jobId) // debug *- cwlogs

	job := queue.Job{
		ID: jobId,
		Run: updateAnalytics(id, ip),
	}
	jobQueue.SendToQueue(job)

	return c.Redirect(http.StatusMovedPermanently, url) 
}


func updateAnalytics(id, ip string) func() {
	return func() {
		data, err := model.GetIndex(id)
		if err != nil {
			return
		}

		data.Clicks++
		if data.CountryMetrics == nil {
			data.CountryMetrics = make(map[string]int64)
		}


		res := map[string]interface{}{}
		req, err := http.Get(fmt.Sprintf("https://api.country.is/%s", ip))
		if req != nil {
			defer req.Body.Close()
			json.NewDecoder(req.Body).Decode(&res)
		}

		if country, ok := res["country"]; ok {
			data.CountryMetrics[country.(string)]++
		} else {
			data.CountryMetrics["UNKNOWN"]++
		}

		err = model.UpdateIndex(id, data)
		if err != nil {
			log.Println(err)
		}
	}
}

func HandleResponseJSON(c echo.Context, code int, msg string, data interface{}) error {
	res := JSONResponse[interface{}]{
		Status: code,
		Message: msg,
	}
	if data != nil {
		res.Data = data
	}

	return c.JSON(code, res)
}

