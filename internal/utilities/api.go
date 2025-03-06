package utilities

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/harrydayexe/Omni/internal/auth"
	"github.com/harrydayexe/Omni/internal/snowflake"
)

// extract the id parameter from the http request
// if the parameter cannot be parsed then an error is written to the http response (if the response is not nil)
func ExtractIdParam(r *http.Request, w http.ResponseWriter, logger *slog.Logger) (snowflake.Snowflake, error) {
	idString := r.PathValue("id")
	logger.InfoContext(r.Context(), "extracting id from request", slog.String("idString", idString))
	idInt, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		logger.InfoContext(r.Context(), "failed to parse id to int", slog.Any("error", err))
		errorMessage := "Url parameter could not be parsed properly."
		if w != nil {
			http.Error(w, errorMessage, http.StatusBadRequest)
		}
		return snowflake.ParseId(0), fmt.Errorf("failed to parse id to int: %w", err)
	}

	return snowflake.ParseId(idInt), nil
}

// check if an error is present and handle the http response if it is
func IsDbError(ctx context.Context, logger *slog.Logger, w http.ResponseWriter, id snowflake.Snowflake, err error) bool {
	if errors.Is(err, sql.ErrNoRows) {
		logger.InfoContext(ctx, "entity not found", slog.Any("id", id))
		http.Error(w, "entity not found", http.StatusNotFound)
		return true
	}
	if err != nil {
		logger.ErrorContext(ctx, "failed to read entity from db", slog.Any("error", err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
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

// check that the content type header is application/json
// adapted from https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func CheckContentTypeHeader(ctx context.Context, logger *slog.Logger, w http.ResponseWriter, r *http.Request) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			logger.InfoContext(ctx, msg)
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return fmt.Errorf("Content-Type header is not application/json")
		}
	}
	return nil
}

// decode the json body of an http request into a struct
// adapted from https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func DecodeJsonBody[T any](ctx context.Context, logger *slog.Logger, w http.ResponseWriter, r *http.Request, obj T) error {
	err := CheckContentTypeHeader(ctx, logger, w, r)
	if err != nil {
		return err
	}

	// Set maximum body size to 1MB to prevent dos attacks
	http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&obj)
	if err != nil {
		// Handle JSON decode error
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
			return fmt.Errorf("Request body contains badly-formed JSON: %w", err)
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			http.Error(w, msg, http.StatusBadRequest)
			return fmt.Errorf("Request body contains badly-formed JSON: %w", err)
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)
			return fmt.Errorf("Request body contains an invalid value for the %q field: %w", unmarshalTypeError.Field, err)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)
			return fmt.Errorf("Request body contains unknown field %s: %w", fieldName, err)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
			return fmt.Errorf("Request body must not be empty: %w", err)
		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
			return fmt.Errorf("Request body must not be larger than 1MB: %w", err)
		default:
			msg := "An unknown error occurred while decoding the request body"
			http.Error(w, msg, http.StatusInternalServerError)
			return err
		}
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return fmt.Errorf("Request body must only contain a single JSON object")
	}

	return nil
}

// CheckBearerAuth checks the Authorization header of an http request and
// verifies the token with the auth service against the given id
func CheckBearerAuth(id snowflake.Snowflake, authService auth.Authable, logger *slog.Logger, w http.ResponseWriter, r *http.Request) error {
	logger.DebugContext(r.Context(), "checking bearer auth")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		msg := "Authorization header is missing"
		logger.InfoContext(r.Context(), msg)
		http.Error(w, msg, http.StatusUnauthorized)
		return auth.ErrUnauthorized
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		msg := "Authorization header format must be Bearer {token}"
		logger.InfoContext(r.Context(), msg)
		http.Error(w, msg, http.StatusUnauthorized)
		return auth.ErrUnauthorized
	}

	err := authService.VerifyToken(r.Context(), authHeaderParts[1], id)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return err
	}

	return nil
}
