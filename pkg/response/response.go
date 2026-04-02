// response code

package response

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]any

// getting data into json format

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, status int, data any) {
	JSON(w, status, envelope{"success": true, "data": data})
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, envelope{"success": false, "error": message})
}

// function for validation  error

func ValidationErro(w http.ResponseWriter, errors map[string]string) {
	JSON(w, http.StatusUnprocessableEntity, envelope{
		"success":           false,
		"error":             "validation failed",
		"validation_errors": errors,
	})
}

func Created(w http.ResponseWriter, data any) {
	Success(w, http.StatusCreated, data)
}

func OK(w http.ResponseWriter, data any) {
	Success(w, http.StatusOK, data)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func NotFound(w http.ResponseWriter, resource string) {
	Error(w, http.StatusNotFound, resource+" not found")
}

func Unauthorized(w http.ResponseWriter) {
	Error(w, http.StatusUnauthorized, "unauthorized: valid token required")
}

func Forbidden(w http.ResponseWriter) {
	Error(w, http.StatusForbidden, "forbidden: insufficient permissions")
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "internal server error")
}

func BadRequest(w http.ResponseWriter, msg string) {
	Error(w, http.StatusBadRequest, msg)
}
