package main

import (
	"encoding/json"
	"net/http"
)

type apiGetter func(url string) (ApiResponse, error)

func apiGet(url string) (ApiResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return ApiResponse{}, nil
	}
	defer resp.Body.Close()
	ar := ApiResponse{}
	jdec := json.NewDecoder(resp.Body)
	err = jdec.Decode(&ar)
	if err != nil {
		return ApiResponse{}, err
	}
	return ar, nil
}
