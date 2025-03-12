package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/harrydayexe/Omni/internal/config"
	"github.com/harrydayexe/Omni/internal/middleware"
	"github.com/harrydayexe/Omni/internal/omniview/connector"
	"github.com/harrydayexe/Omni/internal/omniview/templates"
	"github.com/oxtoacart/bpool"
)

func addRoutes(
	mux *http.ServeMux,
	templates *templates.Templates,
	logger *slog.Logger,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	cfg config.ViewConfig,
) {
	stack := middleware.CreateStack(
		middleware.NewLoggingMiddleware(logger),
		middleware.NewJwtSecret(cfg.JWTSecret),
		middleware.NewMaxBytesReader(),
		newIsLoggedInMiddleware(logger),
	)

	mux.Handle("GET /", stack(handleGetIndex(templates, dataConnector, bufpool, logger)))
	mux.Handle("GET /user/{id}", stack(handleGetUser(templates, dataConnector, bufpool, logger)))
	mux.Handle("GET /post/new", stack(handleGetCreatePost(templates, bufpool, logger)))
	mux.Handle("POST /post/new", stack(handlePostCreatePost(templates, dataConnector, bufpool, logger)))
	mux.Handle("GET /post/{id}", stack(handleGetPost(templates, dataConnector, bufpool, logger)))
	mux.Handle("GET /post/{id}/comments", stack(handleGetComments(templates, dataConnector, bufpool, logger)))
	mux.Handle("GET /login", stack(handleGetLogin(templates, bufpool, logger)))
	mux.Handle("POST /login", stack(handlePostLogin(templates, dataConnector, bufpool, logger)))
	mux.Handle("DELETE /logout", stack(handleDeleteLogout(logger)))
	mux.Handle("GET /signup", stack(handleGetSignup(templates, bufpool, logger)))
	mux.Handle("POST /signup", stack(handlePostSignup(templates, dataConnector, bufpool, logger)))
}

func handleGetIndex(
	templates *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleGetIndexPage(templates, dataConnector, bufpool, logger, isHTMXRequest(r)).ServeHTTP(w, r)
	})
}

func handleGetUser(
	templates *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHTMXRequest(r) {
			// TODO: Handle HTMX request
			w.WriteHeader(http.StatusNotFound)
		} else {
			handleGetUserPage(templates, dataConnector, bufpool, logger).ServeHTTP(w, r)
		}
	})
}

func handleGetPost(
	templates *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHTMXRequest(r) {
			// TODO: Handle HTMX request
			w.WriteHeader(http.StatusNotFound)
		} else {
			handleGetPostPage(templates, dataConnector, bufpool, logger).ServeHTTP(w, r)
		}
	})
}

func handleGetComments(
	t *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isHTMXRequest(r) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		handleGetCommentsPartial(t, dataConnector, bufpool, logger).ServeHTTP(w, r)
	})
}

func handleGetLogin(
	templates *templates.Templates,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHTMXRequest(r) {
			// TODO: Handle HTMX request
			w.WriteHeader(http.StatusNotFound)
		} else {
			handleGetLoginPage(templates, bufpool, logger).ServeHTTP(w, r)
		}
	})
}

func handlePostLogin(
	templates *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlePostLoginPartial(
			templates,
			dataConnector,
			bufpool,
			logger,
			isHTMXRequest(r),
		).ServeHTTP(w, r)
	})
}

func handleDeleteLogout(
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "DELETE request received for /logout")

		cookie := http.Cookie{
			Name:     authCookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.UnixMicro(0),
			HttpOnly: true,
			Secure:   false, // NOTE: Set to true in production when using HTTPS
		}
		http.SetCookie(w, &cookie)
		w.Header().Add("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	})
}

func handleGetCreatePost(
	templates *templates.Templates,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleGetCreatePostPage(templates, bufpool, logger).ServeHTTP(w, r)
	})
}

func handlePostCreatePost(
	templates *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlePostCreatePostPartial(templates, dataConnector, bufpool, logger, isHTMXRequest(r)).ServeHTTP(w, r)
	})
}

func handleGetSignup(
	templates *templates.Templates,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleGetSignupPage(templates, bufpool, logger).ServeHTTP(w, r)
	})
}

func handlePostSignup(
	templates *templates.Templates,
	dataConnector connector.Connector,
	bufpool *bpool.BufferPool,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlePostSignupPartial(
			templates,
			dataConnector,
			bufpool,
			logger,
			isHTMXRequest(r),
		).ServeHTTP(w, r)
	})
}
