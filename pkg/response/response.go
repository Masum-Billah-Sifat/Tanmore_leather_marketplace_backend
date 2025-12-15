package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ---------------------- SUCCESS ----------------------

func OK(w http.ResponseWriter, message string, data interface{}) {
	writeJSON(w, http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	writeJSON(w, http.StatusCreated, SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// ----------------------- ERRORS -----------------------

func BadRequest(w http.ResponseWriter, err error) {
	writeError(w, http.StatusBadRequest, "bad request", err)
}

func Unauthorized(w http.ResponseWriter, err error) {
	writeError(w, http.StatusUnauthorized, "unauthorized", err)
}

func Forbidden(w http.ResponseWriter, err error) {
	writeError(w, http.StatusForbidden, "forbidden", err)
}

func NotFound(w http.ResponseWriter, err error) {
	writeError(w, http.StatusNotFound, "not found", err)
}

func Conflict(w http.ResponseWriter, err error) {
	writeError(w, http.StatusConflict, "conflict", err)
}

func TooManyRequests(w http.ResponseWriter, err error) {
	writeError(w, http.StatusTooManyRequests, "too many requests", err)
}

func ServerError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusInternalServerError, "internal server error", err)
}

// -------------------- INTERNAL HELPERS --------------------

func writeError(w http.ResponseWriter, statusCode int, message string, err error) {
	writeJSON(w, statusCode, ErrorResponse{
		Status:  "error",
		Message: message,
		Error:   err.Error(),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
