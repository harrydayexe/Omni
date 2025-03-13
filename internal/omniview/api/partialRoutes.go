package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	readdatamodels "github.com/harrydayexe/Omni/internal/omniread/datamodels"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	datamodels "github.com/harrydayexe/Omni/internal/omniview/data-models"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	writedatamodels "github.com/harrydayexe/Omni/internal/omniwrite/datamodels"
	"github.com/harrydayexe/Omni/internal/storage"
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
			ID      int64
			Message string
		}{
			ID:      resp.ID,
			Message: "Post created successfully",
		}
		writeTemplateWithBuffer(r.Context(), logger, 0, "newpost-success", t, bufpool, w, successContent)

	})
}

func handlePostPostEditPartial(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "POST request received for partial /post/{id}")

		// Check that the post request has the correct content-type
		err := checkContentTypeHeader(logger, r, formUrlEncoded)
		if err != nil {
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}

		// Parse post id
		postSnowflake, err := utilities.ExtractIdParam(r, nil, logger)
		if err != nil {
			return
		}

		// Create return model
		content := datamodels.NewForm()

		// Parse form
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
				http.StatusUnprocessableEntity, "Edit Post", true,
				t, bufpool, w, content,
			)
			return
		}

		updatedPost := writedatamodels.UpdatedPost{
			Title:       title[0],
			Description: description[0],
			MarkdownUrl: url[0],
		}

		resp, err := dataConnector.UpdatePost(r.Context(), postSnowflake, updatedPost)
		var ae *connector.APIError
		if errors.As(err, &ae) {
			logger.DebugContext(r.Context(), "API error occurred while updating post", slog.String("error", ae.Error()))
			content.Errors["General"] = "An error occurred while creating the post. Please try again later."
			writeFormWithErrors(
				r.Context(), logger,
				http.StatusUnprocessableEntity, "New Post", true,
				t, bufpool, w, content,
			)
			return
		}

		// Write the new post form with or without the errors
		writeFormWithErrors(
			r.Context(), logger,
			http.StatusOK, "New Post", true,
			t, bufpool, w, content,
		)
		successContent := struct {
			ID      int64
			Message string
		}{
			ID:      resp.ID,
			Message: "Post updated successfully",
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
			// TODO: Probably a better way to handle this error
			logger.ErrorContext(r.Context(), "Error occurred while deleting comment", slog.String("error", err.Error()))
			writeTemplateWithBuffer(r.Context(), logger, http.StatusInternalServerError, "errorpage.html", t, bufpool, w, nil)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func handleInsertCommentPartial(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "POST request received for partial /post/{id}/comment")

		// Check that the post request has the correct content-type
		err := checkContentTypeHeader(logger, r, formUrlEncoded)
		if err != nil {
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}

		// Get the post id from the path
		postSnowflake, err := utilities.ExtractIdParam(r, nil, logger)
		if err != nil {
			return
		}
		logger.DebugContext(r.Context(), "Post ID extracted", slog.Int64("id", int64(postSnowflake.ToInt())))

		// Create storage for results
		var username string
		var comment storage.Comment

		// Decode the form
		err = r.ParseForm()
		if err != nil {
			logger.ErrorContext(r.Context(), "Error occurred while parsing form", slog.String("error", err.Error()))
			writeTemplateWithBuffer(r.Context(), logger, http.StatusInternalServerError, "errorpage.html", t, bufpool, w, nil)
			return
		}
		logger.DebugContext(r.Context(), "Form parsed")

		commentContent := r.FormValue("comment")
		if commentContent == "" {
			logger.InfoContext(r.Context(), "Comment content is empty")
			http.Error(w, "Comment content is empty", http.StatusBadRequest)
			return
		}
		logger.DebugContext(r.Context(), "Comment content extracted", slog.String("content", commentContent))

		loggedInUserId := GetUserIdFromCtx(r.Context())
		if loggedInUserId == nil {
			logger.InfoContext(r.Context(), "User not logged in")
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}
		logger.DebugContext(r.Context(), "User ID extracted", slog.Int64("id", int64(loggedInUserId.Id().ToInt())))

		// Create error channels
		const numGoroutines int = 2
		errChan := make(chan error, numGoroutines)
		// Create sub-context
		subctx, cancel := context.WithCancel(r.Context())

		// Get the username
		go func() {
			user, err := dataConnector.GetUser(subctx, loggedInUserId)
			if err != nil {
				logger.InfoContext(r.Context(), "Error occurred while fetching user data", slog.String("error", err.Error()))
				// Cancel other routines
				cancel()
				errChan <- err
				return
			}
			logger.DebugContext(r.Context(), "User data fetched", slog.Int64("id", int64(loggedInUserId.Id().ToInt())))
			username = user.Username
			errChan <- nil
		}()

		go func() {
			// Insert the comment
			newComment := writedatamodels.NewComment{
				UserID:    int64(loggedInUserId.Id().ToInt()),
				Content:   commentContent,
				CreatedAt: time.Now(),
			}
			commentResp, err := dataConnector.InsertComment(r.Context(), postSnowflake, newComment)
			if err != nil {
				logger.InfoContext(r.Context(), "Error occurred while inserting comment", slog.String("error", err.Error()))
				// Cancel other routines
				cancel()
				errChan <- err
				return
			}
			logger.DebugContext(r.Context(), "Comment inserted", slog.Int64("id", commentResp.ID))
			comment = commentResp
			errChan <- nil
		}()

		// Wait for goroutines to finish
		var firstErr error
		for i := 0; i < numGoroutines; i++ {
			// If the error channel has an error and it is the first error...
			if err := <-errChan; err != nil && firstErr == nil {
				logger.InfoContext(r.Context(), "An error occurred while inserting the comment", slog.String("error", err.Error()))
				firstErr = err
			}
		}
		logger.DebugContext(r.Context(), "DB Calls completed", slog.String("Handler", "handleGetUser"))

		var userErr error
		comments := readdatamodels.CommentsForPostReturn{
			TotalPages: 0,
			Comments:   make([]readdatamodels.CommentReturn, 0),
		}
		// Handle firstErr
		if firstErr != nil {
			userErr = errors.New("An error occurred while inserting the comment")
		} else {
			// Map comment into comment return
			comments.Comments = append(comments.Comments, readdatamodels.CommentReturn{
				ID:        comment.ID,
				PostID:    comment.PostID,
				UserID:    comment.UserID,
				Username:  username,
				CreatedAt: comment.CreatedAt,
				Content:   comment.Content,
			})
		}

		// Create comment response model
		content := datamodels.NewCommentsModel(
			userErr, loggedInUserId, comments, int64(postSnowflake.ToInt()), 0,
		)

		writeTemplateWithBuffer(r.Context(), logger, http.StatusOK, "comment-list", t, bufpool, w, content)
	})
}

func handleDeletePostPartial(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "DELETE request received for partial /post/{id}")

		// Parse post id
		postSnowflake, err := utilities.ExtractIdParam(r, nil, logger)
		if err != nil {
			return
		}

		// Request Delete Post
		err = dataConnector.DeletePost(r.Context(), postSnowflake)
		if err != nil {
			// TODO: Probably a better way to handle this error
			logger.ErrorContext(r.Context(), "Error occurred while deleting post", slog.String("error", err.Error()))
			writeTemplateWithBuffer(r.Context(), logger, http.StatusInternalServerError, "errorpage.html", t, bufpool, w, nil)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
