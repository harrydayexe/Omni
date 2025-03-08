package api

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/harrydayexe/Omni/internal/omniview/connector"
	datamodels "github.com/harrydayexe/Omni/internal/omniview/data-models"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	"github.com/oxtoacart/bpool"
)

func handlePostLoginPartial(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for partial /login")

		// Check that the post request has the correct content-type
		err := checkContentTypeHeader(logger, r, formUrlEncoded)
		if err != nil {
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}

		content := datamodels.NewLoginForm()

		r.ParseForm()
		isErr := false
		username, uprs := r.Form["username"]
		if !uprs || len(username) == 0 || len(username[0]) == 0 {
			logger.InfoContext(r.Context(), "username is empty")
			content.Errors["Username"] = "Username is required"
			isErr = true
		} else {
			content.Values["Username"] = username[0]
		}
		password, pprs := r.Form["password"]
		if !pprs || len(password) == 0 || len(password[0]) < 7 {
			logger.InfoContext(r.Context(), "password is not at least 8 characters")
			content.Errors["Password"] = "Password must be at least 8 characters"
			isErr = true
		}

		if isErr {
			writeTemplateWithBuffer(r.Context(), logger, "login-form", t, bufpool, w, content)
			return
		}

		resp, err := dataConnector.Login(r.Context(), username[0], password[0])
		var ae *connector.APIError
		if errors.As(err, &ae) {
			if ae.StatusCode == http.StatusUnauthorized {
				content.Errors["Login"] = "Invalid username or password"
			} else {
				content.Errors["Login"] = "An error occurred while logging in. Please try again later."
			}
		}
		// Write the login form with or without the errors
		writeTemplateWithBuffer(r.Context(), logger, "login-form", t, bufpool, w, content)

		cookie := http.Cookie{
			Name:     "auth_token",
			Value:    resp.Token,
			Path:     "/",
			Expires:  time.Now().Add(time.Duration(resp.Expires * int(time.Second))),
			HttpOnly: true,
			Secure:   false, // NOTE: Set to true in production when using HTTPS
		}
		http.SetCookie(w, &cookie)
		writeTemplateWithBuffer(r.Context(), logger, "login-success", t, bufpool, w, nil)
	})
}
