package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	modellink "github.com/broadcast80/ozon-task/domain/model/link"
	"github.com/broadcast80/ozon-task/internal/pkg/models"
)

type mockShortener struct {
	cutLinkCalled     bool
	getFullLinkCalled bool
	cutLinkInput      string
	getFullLinkInput  string
	cutLinkResult     *modellink.Link
	cutLinkErr        error
	getFullLinkResult *modellink.Link
	getFullLinkErr    error
}

func (m *mockShortener) CutLink(ctx context.Context, url string) (*modellink.Link, error) {
	m.cutLinkCalled = true
	m.cutLinkInput = url
	return m.cutLinkResult, m.cutLinkErr
}

func (m *mockShortener) GetFullLink(ctx context.Context, alias string) (*modellink.Link, error) {
	m.getFullLinkCalled = true
	m.getFullLinkInput = alias
	return m.getFullLinkResult, m.getFullLinkErr
}

func TestHandlers_Create_Success(t *testing.T) {

	mockService := &mockShortener{
		cutLinkResult: &modellink.Link{
			Alias: "somelabuda",
			URL:   "https://primerchik.com",
		},
	}

	router := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := New(router, mockService, logger)
	h.MapHandlers()

	requestBody := models.Request{URL: "https://primerchik.com"}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if !mockService.cutLinkCalled {
		t.Error("CutLink should be called")
	}

	if mockService.cutLinkInput != requestBody.URL {
		t.Errorf("expected CutLink input %s, got %s", requestBody.URL, mockService.cutLinkInput)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	response := modellink.Link{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response.Alias != mockService.cutLinkResult.Alias {
		t.Errorf("expected result %s, got %s", mockService.cutLinkResult.Alias, response.Alias)
	}
}

func TestHandlers_Create_ServiceError(t *testing.T) {

	mockService := &mockShortener{
		cutLinkErr: errors.New("database error"),
	}

	router := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := New(router, mockService, logger)
	h.MapHandlers()

	requestBody := models.Request{URL: "https://example.com"}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestHandlers_Get_Success(t *testing.T) {
	mockService := &mockShortener{
		getFullLinkResult: &modellink.Link{
			Alias: "diehard",
			URL:   "https://newyear.com",
		},
	}

	router := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := New(router, mockService, logger)
	h.MapHandlers()

	requestBody := models.Request{Alias: "short-123"}
	bodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("GET", "/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if !mockService.getFullLinkCalled {
		t.Error("GetFullLink should be called")
	}

	response := modellink.Link{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response.URL != mockService.getFullLinkResult.URL {
		t.Errorf("expected result %s, got %s", mockService.getFullLinkResult.URL, response.URL)
	}
}

func TestHandlers_Get_InvalidJSON(t *testing.T) {
	mockService := &mockShortener{}
	router := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := New(router, mockService, logger)
	h.MapHandlers()

	req, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}
