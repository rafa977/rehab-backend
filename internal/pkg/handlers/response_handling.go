package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rehab-backend/internal/pkg/models"
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

func ProduceSuccessResponse(responseMsg string, w http.ResponseWriter, r *http.Request) {
	var response models.Response

	currentDate := time.Now().Format("2006-01-02 15:04:05")
	response.Date = currentDate

	response.Status = "success"
	response.Message = ""
	response.Response = responseMsg
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}
