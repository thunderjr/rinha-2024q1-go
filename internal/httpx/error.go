package httpx

import (
	"encoding/json"
	"net/http"
)

type HttpError struct {
	cause error
	w     http.ResponseWriter
}

func NewError(w http.ResponseWriter, err error) *HttpError {
	return &HttpError{
		cause: err,
		w:     w,
	}
}

func (h *HttpError) Status(code int) *HttpError {
	h.w.WriteHeader(code)
	return h
}

func (h *HttpError) Send() {
	json.NewEncoder(h.w).Encode(map[string]string{
		"error": h.cause.Error(),
	})
}
