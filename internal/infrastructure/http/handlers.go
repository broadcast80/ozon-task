package infrastracture

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/broadcast80/ozon-task/internal/config"
	"github.com/broadcast80/ozon-task/internal/domain"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Alias string `json:"alias"`
}

type ServiceInterface interface {
	Create(ctx context.Context, link domain.Link) (string, error)
	Get(ctx context.Context, alias string) (string, error)
}

type handlers struct {
	cfg     *config.Config
	router  *http.ServeMux
	service ServiceInterface
}

// тут сделать интерфейс или нет
func New(cfg *config.Config, router *http.ServeMux, service ServiceInterface) *handlers {
	return &handlers{
		cfg:     cfg,
		router:  router,
		service: service,
	}
}

func (h *handlers) ListenAndServe(cfg config.HTTPServer) error {
	address := ":" + cfg.Port
	err := http.ListenAndServe(address, h.router)
	if err != nil {
		return fmt.Errorf("listen and serve error: %w", err)
	}
	return nil
}

func (h *handlers) MapHandlers() error {
	h.router.HandleFunc("POST /", h.Create)

	return nil
}

func (h *handlers) Create(w http.ResponseWriter, r *http.Request) {

	url := r.Form.Get("url")

	link := domain.Link{
		URL:       url,
		CreatedAt: time.Now(),
	}

	alias, err := h.service.Create(context.TODO(), link)
	if err != nil {
		// обработка
	}

	response := Response{
		Alias: alias,
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Json", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	return

}

func main() {
	// router := http.NewServeMux()
}
