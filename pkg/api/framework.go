package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/llkhacquan/cisab/pkg/utils/logger"
)

// HandlerFunc is a custom function type for API handlers that returns data, status code, and optional error
type HandlerFunc func(r *http.Request) (interface{}, error)

// Endpoint represents an API endpoint with a handler.
// This inherits the idea from go-kit, but we don't need to use the full framework for our needs.
type Endpoint struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

// HandleEndpoint creates an http.HandlerFunc from an Endpoint
func HandleEndpoint(e Endpoint, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute the handler
		data, err := e.Handler(r)

		// Handle errors - always return 500 Internal Server Error for now
		// Todo: we should have custom error handling here later
		if err != nil {
			log.Error("error in request handler", "path", r.URL.Path, "method", r.Method, "error", err.Error())
			errorResp := ErrorResponse("Internal server error", err)
			WriteJSON(w, http.StatusInternalServerError, errorResp, log)
			return
		}
		// Else, return success response
		WriteJSON(w, http.StatusOK, SuccessResponse("", data), log)
	}
}

// RegisterEndpoints registers multiple endpoints with router
func RegisterEndpoints(router *mux.Router, endpoints []Endpoint, log *logger.Logger) {
	for _, endpoint := range endpoints {
		handler := HandleEndpoint(endpoint, log)
		router.HandleFunc(endpoint.Path, handler).Methods(endpoint.Method)
	}
}

// GetURLParams retrieves URL parameters from the request
func GetURLParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}
