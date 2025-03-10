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
		logger.InfoContext(r.Context(), "POST request received for partial /login")

		// Check that the post request has the correct content-type
		err := checkContentTypeHeader(logger, r, formUrlEncoded)
		if err != nil {
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}

		content := datamodels.NewForm()

		r.ParseForm()
		isErr := false
		username, uprs := r.Form["username"]
		if !uprs || len(username) == 0 || len(username[0]) == 0 {
			logger.DebugContext(r.Context(), "username is empty")
			content.Errors["Username"] = "Username is required"
			isErr = true
		} else {
			content.Values["Username"] = username[0]
		}
		password, pprs := r.Form["password"]
		if !pprs || len(password) == 0 || len(password[0]) < 7 {
			logger.DebugContext(r.Context(), "password is not at least 8 characters")
			content.Errors["Password"] = "Password must be at least 8 characters"
			isErr = true
		}

		if isErr {
			writeTemplateWithBuffer(r.Context(), logger, http.StatusUnprocessableEntity, "login-form", t, bufpool, w, content)
			return
		}

		resp, err := dataConnector.Login(r.Context(), username[0], password[0])
		var ae *connector.APIError
		if errors.As(err, &ae) {
			logger.DebugContext(r.Context(), "API error occurred while logging in", slog.String("error", ae.Error()))
			if ae.StatusCode == http.StatusUnauthorized {
				content.Errors["Login"] = "Invalid username or password"
			} else if ae.StatusCode == http.StatusNotFound {
				content.Errors["Username"] = "User not found"
			} else {
				content.Errors["Login"] = "An error occurred while logging in. Please try again later."
			}
			writeTemplateWithBuffer(r.Context(), logger, http.StatusUnprocessableEntity, "login-form", t, bufpool, w, content)
			return
		}
		cookie := http.Cookie{
			Name:     authCookieName,
			Value:    resp.Token,
			Path:     "/",
			Expires:  time.Now().Add(time.Duration(resp.Expires * int(time.Second))),
			HttpOnly: true,
			Secure:   false, // NOTE: Set to true in production when using HTTPS
		}
		http.SetCookie(w, &cookie)

		// Write the login form with or without the errors
		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "login-form", t, bufpool, w, content)
		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "login-success", t, bufpool, w, nil)
	})
}

func handlePostCreatePostPartial(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
	isHTMXRequest bool,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "POST request received for partial /post/new")

		if _, prs := hasValidAuthToken(r, logger); !prs {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Check that the post request has the correct content-type
		err := checkContentTypeHeader(logger, r, formUrlEncoded)
		if err != nil {
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}

		content := datamodels.NewForm()

		r.ParseForm()
		isErr := false
		title, prs1 := r.Form["title"]
		if !prs1 || len(title) == 0 || len(title[0]) == 0 {
			logger.DebugContext(r.Context(), "title is empty")
			content.Errors["Title"] = "Title is required"
			isErr = true
		} else {
			content.Values["Title"] = title[0]
		}
		description, prs2 := r.Form["description"]
		if !prs2 || len(description) == 0 || len(description[0]) == 0 {
			logger.DebugContext(r.Context(), "description is empty")
			content.Errors["Description"] = "Description is required"
			isErr = true
		} else {
			content.Values["Description"] = description[0]
		}
		url, prs3 := r.Form["url"]
		if !prs3 || len(url) == 0 || len(url[0]) == 0 {
			logger.DebugContext(r.Context(), "url is empty")
			content.Errors["URL"] = "URL is required"
			isErr = true
		} else {
			content.Values["URL"] = url[0]
		}

		if isErr {
			if isHTMXRequest {
				writeTemplateWithBuffer(r.Context(), logger, http.StatusUnprocessableEntity, "newpostform", t, bufpool, w, content)
			} else {
				pageContent := datamodels.NewFormPage(r.Context(), "New Post")
				pageContent.Form = content
				if _, prs := hasValidAuthToken(r, logger); prs {
					pageContent.NavBar.IsLoggedIn = true
				}
				writeTemplateWithBuffer(r.Context(), logger, http.StatusUnprocessableEntity, "newpost.html", t, bufpool, w, pageContent)
			}
			return
		}

		// resp, err := dataConnector.Login(r.Context(), username[0], password[0])
		// var ae *connector.APIError
		// if errors.As(err, &ae) {
		// 	logger.DebugContext(r.Context(), "API error occurred while logging in", slog.String("error", ae.Error()))
		// 	if ae.StatusCode == http.StatusUnauthorized {
		// 		content.Errors["Login"] = "Invalid username or password"
		// 	} else if ae.StatusCode == http.StatusNotFound {
		// 		content.Errors["Username"] = "User not found"
		// 	} else {
		// 		content.Errors["Login"] = "An error occurred while logging in. Please try again later."
		// 	}
		// 	writeTemplateWithBuffer(r.Context(), logger, http.StatusUnprocessableEntity, "login-form", t, bufpool, w, content)
		// 	return
		// }
		// cookie := http.Cookie{
		// 	Name:     authCookieName,
		// 	Value:    resp.Token,
		// 	Path:     "/",
		// 	Expires:  time.Now().Add(time.Duration(resp.Expires * int(time.Second))),
		// 	HttpOnly: true,
		// 	Secure:   false, // NOTE: Set to true in production when using HTTPS
		// }
		// http.SetCookie(w, &cookie)
		//
		// // Write the login form with or without the errors
		// writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "login-form", t, bufpool, w, content)
		// writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "login-success", t, bufpool, w, nil)

	})
}
