package infrastracture

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/broadcast80/ozon-task/domain/link"
	"github.com/broadcast80/ozon-task/internal/pkg/models"
)

type handlers struct {
	router  *http.ServeMux
	service *link.Shortener
	logger  *slog.Logger
}

func New(router *http.ServeMux, service *link.Shortener, logger *slog.Logger) *handlers {
	return &handlers{
		router:  router,
		service: service,
		logger:  logger,
	}
}

func (h *handlers) ListenAndServe(port string) error {
	address := ":" + port
	err := http.ListenAndServe(address, h.router)
	if err != nil {
		return fmt.Errorf("listen and serve error: %w", err)
	}
	return nil
}

func (h *handlers) MapHandlers() error {
	h.router.HandleFunc("POST /", h.Create)
	h.router.HandleFunc("GET /", h.Get)

	return nil
}

func (h *handlers) Create(w http.ResponseWriter, r *http.Request) {

	// ...
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "1 service fault", http.StatusInternalServerError)
	}

	var request models.Request

	err = json.Unmarshal(body, &request)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "2 service fault", http.StatusInternalServerError)
	}
	// ...

	link, err := h.service.CutLink(r.Context(), request.URL)
	if err != nil {
		http.Error(w, "service fault", http.StatusInternalServerError)
	}

	data, err := json.Marshal(link)
	if err != nil {
		http.Error(w, "Json", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	return
}

func (h *handlers) Get(w http.ResponseWriter, r *http.Request) {

	// ...
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "1 service fault", http.StatusInternalServerError)
	}

	var request models.Request

	err = json.Unmarshal(body, &request)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "2 service fault", http.StatusInternalServerError)
	}
	// ...

	link, err := h.service.GetFullLink(r.Context(), request.Alias)
	if err != nil {
		http.Error(w, "service fault", http.StatusInternalServerError)
	}

	data, err := json.Marshal(link)
	if err != nil {
		http.Error(w, "Json", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	return
}
