package utilities

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/harrydayexe/Omni/internal/snowflake"
)

// extract the id parameter from the http request
// if the parameter cannot be parsed then an error is written to the http response
func ExtractIdParam(r *http.Request, w http.ResponseWriter, logger *slog.Logger) (snowflake.Snowflake, error) {
	idString := r.PathValue("id")
	logger.InfoContext(r.Context(), "extracting id from request", slog.String("idString", idString))
	idInt, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		logger.ErrorContext(r.Context(), "failed to parse id to int", slog.Any("error", err))
		errorMessage := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMessage))
		return snowflake.ParseId(0), fmt.Errorf("failed to parse id to int: %w", err)
	}

	return snowflake.ParseId(idInt), nil
}

// check if an error is present and handle the http response if it is
func IsDbError(ctx context.Context, logger *slog.Logger, w http.ResponseWriter, id snowflake.Snowflake, err error) bool {
	if errors.Is(err, sql.ErrNoRows) {
		logger.ErrorContext(ctx, "entity not found", slog.Any("id", id))
		w.WriteHeader(http.StatusNotFound)
		return true
	}
	if err != nil {
		logger.ErrorContext(ctx, "failed to read entity from db", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return true
	}

	return false
}

// marshall an entity to a json object and write to http response
func MarshallToResponse(ctx context.Context, logger *slog.Logger, w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		logger.ErrorContext(ctx, "failed to serialize entity to json", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}

func CheckContentTypeHeader(ctx context.Context, logger *slog.Logger, w http.ResponseWriter, r *http.Request) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			logger.ErrorContext(ctx, msg)
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return fmt.Errorf("Content-Type header is not application/json")
		}
	}
	return nil
}
