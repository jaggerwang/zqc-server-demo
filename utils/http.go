package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

func HttpRequestJson(req *http.Request) (result map[string]interface{}, err error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"method":     req.Method,
		"url":        req.URL,
		"reqHeader":  req.Header,
		"respHeader": resp.Header,
		"respBody":   string(body),
	}).Debug("http request")

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
