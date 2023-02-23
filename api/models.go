package api

import (
	"afd-support/afdian"
)

type APIService struct {
	Token      string
	AfdianItem *afdian.AfdianAPIService
}

type APIRequest struct {
	Token string `json:"token"`
	Data  string `json:"data"`
	Ts    int    `json:"ts"`
	Auth  string `json:"auth"`
}

type APIReponse struct {
	Status  int
	Message string
	Data    *afdian.AfdianQueryResponse
}
