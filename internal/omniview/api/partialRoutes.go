package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/harrydayexe/Omni/internal/omniview/connector"
	datamodels "github.com/harrydayexe/Omni/internal/omniview/data-models"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	writedatamodels "github.com/harrydayexe/Omni/internal/omniwrite/datamodels"
	"github.com/harrydayexe/Omni/internal/utilities"
	"github.com/oxtoacart/bpool"
)

func handlePostLoginPartial(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
	isHTMXRequest bool,
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
		content.Values["Title"] = "Login"
		content.Values["HXDest"] = "/login"

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
			writeFormWithErrors(
				r.Context(), logger,
				http.StatusUnprocessableEntity, "Login", isHTMXRequest,
				t, bufpool, w, content,
			)
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
			writeFormWithErrors(
				r.Context(), logger,
				http.StatusUnprocessableEntity, "Login", isHTMXRequest,
				t, bufpool, w, content,
			)
			return
		}
		logger.DebugContext(r.Context(), "Login call finished")
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
		writeFormWithErrors(
			r.Context(), logger,
			http.StatusOK, "Login", isHTMXRequest,
			t, bufpool, w, content,
		)
		successContent := datamodels.NewFormSuccess("Login successful", "/")
		writeTemplateWithBuffer(r.Context(), logger, 0, "login-success", t, bufpool, w, successContent)
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
			writeFormWithErrors(
				r.Context(), logger,
				http.StatusUnprocessableEntity, "New Post", isHTMXRequest,
				t, bufpool, w, content,
			)
			return
		}

		newPost := writedatamodels.NewPost{
			UserID:      r.Context().Value(UserIdCtxKey).(uint64),
			CreatedAt:   time.Now(),
			Title:       title[0],
			Description: description[0],
			MarkdownUrl: url[0],
		}

		resp, err := dataConnector.CreatePost(r.Context(), newPost)
		var ae *connector.APIError
		if errors.As(err, &ae) {
			logger.DebugContext(r.Context(), "API error occurred while creating post", slog.String("error", ae.Error()))
			content.Errors["General"] = "An error occurred while creating the post. Please try again later."
			writeFormWithErrors(
				r.Context(), logger,
				http.StatusUnprocessableEntity, "New Post", isHTMXRequest,
				t, bufpool, w, content,
			)
			return
		}

		// Write the new post form with or without the errors
		writeFormWithErrors(
			r.Context(), logger,
			http.StatusOK, "New Post", isHTMXRequest,
			t, bufpool, w, content,
		)
		successContent := struct {
			ID int64
		}{
			ID: resp.ID,
		}
		writeTemplateWithBuffer(r.Context(), logger, 0, "newpost-success", t, bufpool, w, successContent)

	})
}

func handlePostSignupPartial(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
	isHTMXRequest bool,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "POST request received for partial /signup")

		// Check that the post request has the correct content-type
		err := checkContentTypeHeader(logger, r, formUrlEncoded)
		if err != nil {
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}

		content := datamodels.NewForm()
		content.Values["Title"] = "Sign Up"
		content.Values["HXDest"] = "/signup"

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
			writeFormWithErrors(
				r.Context(), logger,
				http.StatusUnprocessableEntity, "Sign Up", isHTMXRequest,
				t, bufpool, w, content,
			)
			return
		}

		// Sign the user up
		resp, err := dataConnector.Signup(r.Context(), username[0], password[0])
		logger.DebugContext(r.Context(), "Signup call finished")
		var ae *connector.APIError
		if errors.As(err, &ae) {
			logger.DebugContext(r.Context(), "API error occurred while signing up", slog.String("error", ae.Error()))
			if ae.StatusCode == http.StatusUnprocessableEntity {
				content.Errors["Login"] = "Invalid username or password"
			} else if ae.StatusCode == http.StatusConflict {
				content.Errors["Username"] = "Username taken"
			} else {
				content.Errors["Login"] = "An error occurred while signing up. Please try again later."
			}
			writeFormWithErrors(
				r.Context(), logger,
				http.StatusUnprocessableEntity, "Sign Up", isHTMXRequest,
				t, bufpool, w, content,
			)
			return
		}

		// Create the data for the success message
		successContent := datamodels.NewFormSuccess(
			"Account created successfully",
			"/user/"+strconv.Itoa(int(resp.ID)),
		)

		// Now login the user
		loginResp, err := dataConnector.Login(r.Context(), username[0], password[0])
		logger.DebugContext(r.Context(), "Login call finished")
		if err != nil {
			// This is not a massively critical error,
			// we should just redirect the user to the login page
			logger.InfoContext(r.Context(), "Error occurred while logging in user after signup", slog.String("error", err.Error()))
			successContent.RedirectURL = "/login"
		} else {
			cookie := http.Cookie{
				Name:     authCookieName,
				Value:    loginResp.Token,
				Path:     "/",
				Expires:  time.Now().Add(time.Duration(loginResp.Expires * int(time.Second))),
				HttpOnly: true,
				Secure:   false, // NOTE: Set to true in production when using HTTPS
			}
			http.SetCookie(w, &cookie)
		}

		// Write the login form with or without the errors
		writeFormWithErrors(
			r.Context(), logger,
			http.StatusOK, "Sign Up", isHTMXRequest,
			t, bufpool, w, content,
		)
		writeTemplateWithBuffer(r.Context(), logger, 0, "login-success", t, bufpool, w, successContent)
		logger.DebugContext(r.Context(), "Finished writing signup response")
	})
}

func handleGetCommentsPartial(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "GET request received for partial /post/{id}/comments")

		// Parse user id
		postSnowflake, err := utilities.ExtractIdParam(r, nil, logger)
		if err != nil {
			content := datamodels.NewErrorPageModel(
				"Post not found",
				"The post you are looking for does not exist.",
			)
			writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "errorpage.html", t, bufpool, w, content)
			return
		}

		// Attempt to get page url query param
		var pageNum int = 1
		pageQuery := r.URL.Query().Get("page")
		if pageQuery != "" {
			logger.DebugContext(r.Context(), "Page number specified", slog.String("pageQueryParam", pageQuery))
			pageNumTemp, err := strconv.Atoi(pageQuery)
			if err == nil {
				logger.DebugContext(r.Context(), "Page number parsed", slog.Int("pageNum", pageNumTemp))
				pageNum = pageNumTemp
			}
		}

		// Get Comments
		commentResp, err := dataConnector.GetPostComments(r.Context(), postSnowflake, pageNum)
		commentModel := datamodels.NewCommentsModel(
			err,
			GetUserIdFromCtx(r.Context()),
			commentResp,
			int64(postSnowflake.ToInt()),
			pageNum+1,
		)
		if err != nil {
			logger.InfoContext(r.Context(), "Error occurred while retrieving comments", slog.String("error", err.Error()))
		} else {
			logger.DebugContext(r.Context(), "Successfully retrieved comments")
		}

		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "comment-list", t, bufpool, w, commentModel)
	})
}

func handleDeleteCommentPartial(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "DELETE request received for partial /comment/{id}")

		// Parse comment id
		commentSnowflake, err := utilities.ExtractIdParam(r, w, logger)
		if err != nil {
			return
		}

		// Request Delete Comment
		err = dataConnector.DeleteComment(r.Context(), commentSnowflake)
		if err != nil {
			logger.ErrorContext(r.Context(), "Error occurred while deleting comment", slog.String("error", err.Error()))
			writeTemplateWithBuffer(r.Context(), logger, http.StatusInternalServerError, "errorpage.html", t, bufpool, w, nil)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
