package api

import (
	"encoding/json"
	"net/http"

	"github.com/llkhacquan/cisab/pkg/utils/logger"
	"github.com/pkg/errors"
)

// Response is a generic response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse creates an error response with given message and error
func ErrorResponse(message string, err error) Response {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	return Response{
		Status:  "error",
		Message: message,
		Error:   errMsg,
	}
}

// SuccessResponse creates a success response with given message and data
func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// WriteJSON writes the provided response as JSON to the http.ResponseWriter
func WriteJSON(w http.ResponseWriter, status int, response interface{}, log *logger.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("failed to encode response to JSON", "error", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// ReadJSON attempts to decode the request body into the provided destination
func ReadJSON(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return errors.Wrap(err, "failed to decode body")
	}
	return nil
}
