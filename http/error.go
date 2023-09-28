package http

import "net/http"

type AppError struct {
    Error       error
    Message     string
    Code        int
    ContentType string
}

func internalServerError(err error) *AppError {
    return &AppError{
        Error:   err,
        Message: "Internal server error",
        Code:    http.StatusInternalServerError,
    }
}

func notFound(err error) *AppError {
    return &AppError{Error: err, Code: http.StatusNotFound}
}

func badRequest(err error) *AppError {
    return &AppError{Error: err, Code: http.StatusBadRequest}
}

func (e *AppError) AsJSON() *AppError {
    e.ContentType = jsonMediaType
    return e
}

func (e *AppError) WithMessage(message string) *AppError {
    e.Message = message
    return e
}

func (e *AppError) IsJSON() bool {
    return e.ContentType == jsonMediaType
}
