package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"rehab/internal/pkg/models"
)

func ProduceErrorResponse(msg string, w http.ResponseWriter, r *http.Request) {
	var response models.Response

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	response.Status = "error"
	response.Message = msg
	response.Response = ""
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}

func ProduceSuccessResponse(responseMsg string, message string, w http.ResponseWriter, r *http.Request) {
	var response models.Response

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	response.Status = "success"
	response.Message = message
	response.Response = responseMsg
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}

func ProduceJsonSuccessResponse(responseMsg interface{}, message string, w http.ResponseWriter, r *http.Request) {
	var response models.ResponseJSON

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	response.Status = "success"
	response.ResponseCode = 200
	response.Message = message
	response.Response = responseMsg
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}
