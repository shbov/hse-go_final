package http_helpers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"strconv"

	"net/http"
)

// WriteResponse - вспомогательная функция, которая записывет http статус-код и текстовое сообщение в ответ клиенту.
// Нужна для уменьшения дублирования кода и улучшения читаемости кода вызывающей функции.
func WriteResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
	_, _ = w.Write([]byte("\n"))
}

// WriteJSONResponse - вспомогательная функция, которая запсывает http статус-код и сообщение в формате json в ответ клиенту.
// Нужна для уменьшения дублирования кода и улучшения читаемости кода вызывающей функции.
func WriteJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't marshal data: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	WriteResponse(w, status, string(response))
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func ParseQueryToFloat(field string, values url.Values) (float64, error) {
	value := values.Get(field)

	if value == "" {
		return 0, fmt.Errorf("Field %s is required!\n", field)
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return floatValue, nil
}
