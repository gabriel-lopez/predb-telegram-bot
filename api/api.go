package api

import (
	"fmt"
	"net/http"
	"net/url"
	"encoding/json"
	"log"
	"errors"
)

var preAPIQuery = "https://predb.ovh/api/v1/?q=%s&count=%d"

func QuerySphinx(client *http.Client, q string, max int) ([]sphinxRow, error) {
	resp, err := client.Get(fmt.Sprintf(preAPIQuery, url.QueryEscape(q), max))
	if err != nil {
		return nil, err
	}

	var api apiResponse

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&api)
	if err != nil {
		return nil, err
	}

	if api.Status != "success" {
		log.Println(resp.Body)
		return nil, errors.New("Internal error")
	}

	return api.Data.Rows, nil
}