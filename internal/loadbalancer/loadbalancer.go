package loadbalancer

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
)

type loadBalancer struct {
	Config  Config
	Logger  *slog.Logger
	Proxies map[string]*LoadBalancerProxy
	Mux     *http.ServeMux
}

func New(config Config, logger *slog.Logger) *http.ServeMux {
	logger.Debug("Creating new load balancer")

	loadBalancer := &loadBalancer{
		Config:  config,
		Logger:  logger,
		Proxies: make(map[string]*LoadBalancerProxy),
	}

	// TODO: Add valid paths from config to serve mux

	router := http.NewServeMux()
	router.HandleFunc("POST /addz", func(w http.ResponseWriter, r *http.Request) {
		loadBalancer.addBackend(w, r)
	})
	router.HandleFunc("DELETE /removez", func(w http.ResponseWriter, r *http.Request) {
		loadBalancer.removeBackend(w, r)
	})
	router.HandleFunc("GET /livez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		loadBalancer.readyz(w, r)
	})

	return router
}

// Readiness check endpoint
// Note this returns StatusServiceUnavailable if no backends are available for any path
func (loadBalancer *loadBalancer) readyz(w http.ResponseWriter, r *http.Request) {
	loadBalancer.Logger.InfoContext(r.Context(), "readyz GET request received")

	if len(loadBalancer.Proxies) == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	for _, proxy := range loadBalancer.Proxies {
		if proxy.balancer.Len() == 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// Add a new backend to the load balancer
func (loadBalancer *loadBalancer) addBackend(w http.ResponseWriter, r *http.Request) {
	loadBalancer.Logger.InfoContext(r.Context(), "addz POST request received")
	var c struct {
		Path    string `json:"path"`
		Address string `json:"address"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&c)
	if err != nil {
		loadBalancer.Logger.ErrorContext(r.Context(), "failed to decode request body", slog.Any("error", err))
		var errorMessage = `{"error":"Bad Request","message":"Request body could not be parsed properly."}`
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMessage))
		return
	}

	address, err := url.Parse(c.Address)
	if err != nil {
		loadBalancer.Logger.ErrorContext(r.Context(), "failed to parse address", slog.Any("error", err))
		var errorMessage = `{"error":"Bad Request","message":"Address could not be parsed properly."}`
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMessage))
		return
	}

	proxy, prs := loadBalancer.Proxies[c.Path]
	if !prs {
		loadBalancer.Logger.ErrorContext(r.Context(), "path not found", slog.String("path", c.Path))
		var errorMessage = `{"error":"Not Found","message":"Path not found."}`
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errorMessage))
		return
	}

	proxy.balancer.Add(address)

	w.WriteHeader(http.StatusCreated)
}

func (loadBalancer *loadBalancer) removeBackend(w http.ResponseWriter, r *http.Request) {
	loadBalancer.Logger.InfoContext(r.Context(), "removez DELETE request received")

	pathString := r.URL.Query().Get("path")
	addressString := r.URL.Query().Get("address")
	if pathString == "" || addressString == "" {
		loadBalancer.Logger.ErrorContext(r.Context(), "missing path query parameter")
		var errorMessage = `{"error":"Bad Request","message":"path query parameter missing."}`
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMessage))
		return
	}

	address, err := url.Parse(addressString)
	if err != nil {
		loadBalancer.Logger.ErrorContext(r.Context(), "failed to parse address", slog.Any("error", err))
		var errorMessage = `{"error":"Bad Request","message":"Address could not be parsed properly."}`
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMessage))
		return
	}

	proxy, prs := loadBalancer.Proxies[pathString]
	if !prs {
		loadBalancer.Logger.ErrorContext(r.Context(), "path not found", slog.String("path", pathString))
		var errorMessage = `{"error":"Not Found","message":"Path not found."}`
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errorMessage))
		return
	}

	proxy.balancer.Remove(address)

	w.WriteHeader(http.StatusOK)
}
