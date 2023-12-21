package httpadapter

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/shbov/hse-go_final/internal/driver/service"
	"net/http"
)

// может, вообще этот файл вынести в pkg?

// writeResponse - вспомогательная функция, которая записывет http статус-код и текстовое сообщение в ответ клиенту.
// Нужна для уменьшения дублирования кода и улучшения читаемости кода вызывающей функции.
func writeResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
	_, _ = w.Write([]byte("\n"))
}

// writeJSONResponse - вспомогательная функция, которая запсывает http статус-код и сообщение в формате json в ответ клиенту.
// Нужна для уменьшения дублирования кода и улучшения читаемости кода вызывающей функции.
func writeJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't marshal data: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	writeResponse(w, status, string(response))
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrRequestIsIncorrect): // вынести ошибку в pkg?
		writeJSONResponse(w, http.StatusBadRequest, Error{Message: err.Error()})
	default:
		writeJSONResponse(w, http.StatusInternalServerError, Error{Message: err.Error()})
	}
}
