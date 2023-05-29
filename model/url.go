package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/teris-io/shortid"
)

type URL struct {
	URL string `redis:"url" json:"url"`
	Clicks int64 `redis:"clicks" json:"clicks"`
	Metrics string `redis:"metrics" json:"-"` // stringified (map[string]int64)
	CountryMetrics map[string]int64 `redis:"-" json:"country_metrics"` 
}

func URLIndex(url string) (string, error) {
	id, err := shortid.Generate()
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("url:%s", id)
	values := URL{
		URL: url,
		Clicks: 0,
		Metrics: "{}",
	}

	if err = rdb.HSet(context.TODO(), key, values).Err(); err != nil {
		return "", err
	}

	return id, nil
}

func URLFromIndex(id string) (string, error) {
	key := fmt.Sprintf("url:%s", id)
	url, err := rdb.HGet(context.TODO(), key, "url").Result()
	if err != nil {
		return "", err
	}

	if url == "" {
		return "", errors.New("index not found")
	}

	return url, nil
}

func UnpublishIndex(id string) error {
	key := fmt.Sprintf("url:%s", id)
	if err := rdb.HDel(context.TODO(), key, "url").Err(); err != nil {
		return err
	}

	return nil
}

func GetIndex(id string) (*URL, error) {
	key := fmt.Sprintf("url:%s", id)
	info, err := rdb.HGetAll(context.TODO(), key).Result()
	if err != nil {
		return nil, err
	}

	var url URL
	if longUrl, ok := info["url"]; ok {
		url.URL = longUrl
	}
	if clicks, ok := info["clicks"]; ok {
		if c, err := strconv.Atoi(clicks); err == nil {
			url.Clicks = int64(c)
		}
	}
	if metrics, ok := info["metrics"]; ok {
		json.Unmarshal([]byte(metrics), &url.CountryMetrics)
	}

	return &url, nil
}

func UpdateIndex(id string, data *URL) error {
	key := fmt.Sprintf("url:%s", id)

	if buff, err := json.Marshal(data.CountryMetrics); err == nil {
		data.Metrics = string(buff)
	}

	if err := rdb.HSet(context.TODO(), key, data).Err(); err != nil {
		return err
	}

	return nil
}

