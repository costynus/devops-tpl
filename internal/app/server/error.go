package server

import (
	"devops-tpl/internal/usecase"
	"errors"
	"net/http"
)

func errorHandler(w http.ResponseWriter, err error) {
	if errors.Is(err, usecase.ErrNotImplemented) {
		http.Error(w, err.Error(), http.StatusNotImplemented)
	} else if errors.Is(err, usecase.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else if errors.Is(err, usecase.ErrSignNotEqual) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
