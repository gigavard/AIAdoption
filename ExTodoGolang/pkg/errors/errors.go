package errors

import (
	"encoding/json"
	"net/http"
)

// ProblemDetail implements RFC 7807 Problem Details for HTTP APIs
type ProblemDetail struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// AppError represents an application error
type AppError struct {
	ProblemDetail
}

func NewError(errorType, detail string, status int) *AppError {
	return &AppError{
		ProblemDetail: ProblemDetail{
			Type:   errorType,
			Title:  errorType,
			Status: status,
			Detail: detail,
		},
	}
}

func (e *AppError) Error() string {
	return e.Detail
}

func Respond(w http.ResponseWriter, err *AppError) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err.ProblemDetail)
}
