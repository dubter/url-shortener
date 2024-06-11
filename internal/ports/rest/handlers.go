package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"url-shortner/internal/domain"
	"url-shortner/internal/ports/rest/request"
	"url-shortner/internal/ports/rest/response"
)

type ServiceURLShortener interface {
	Proxy(ctx context.Context, id int) (string, error)
	Create(ctx context.Context, url string) (*domain.Link, error)
}

type ServiceEncoder interface {
	Encode(int) string
	Decode(string) int
}

type ServiceRender interface {
	Home(http.ResponseWriter)
	Icon(http.ResponseWriter, *http.Request)
}

type Handler struct {
	logger       *slog.Logger
	urlshortener ServiceURLShortener
	encoder      ServiceEncoder
	render       ServiceRender
}

func NewHandler(logger *slog.Logger, urlshortener ServiceURLShortener, encoder ServiceEncoder, render ServiceRender) *Handler {
	return &Handler{
		logger:       logger,
		urlshortener: urlshortener,
		encoder:      encoder,
		render:       render,
	}
}

func (h *Handler) Homepage(w http.ResponseWriter, _ *http.Request) {
	h.render.Home(w)
}

func (h *Handler) Icon(w http.ResponseWriter, r *http.Request) {
	h.render.Icon(w, r)
}

func (h *Handler) RegisterURL(w http.ResponseWriter, r *http.Request) {
	url, err := getUrlFromPayload(r)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.Body{"message": err.Error()})
		return
	}

	// check if link already exists on database
	newLink, err := h.urlshortener.Create(r.Context(), url)
	if err != nil {
		h.logger.Error("failed to create short url", slog.String("error", err.Error()))
		response.JSON(w, http.StatusInternalServerError, response.Body{"message": "failed to create short url"})
		return
	}

	shortCode, shortURL := buildShortURL(h.encoder, r.Host, newLink)
	body := response.Body{"short_code": shortCode, "short_url": shortURL}
	response.JSON(w, http.StatusOK, body)
}

func (h *Handler) ProxyURLCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if len(code) == 0 {
		response.JSON(w, http.StatusNotFound, response.Body{"message": "url param 'code' should not be empty"})
		return
	}

	id := h.encoder.Decode(code)

	link, err := h.urlshortener.Proxy(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrURLNotFound) {
			response.JSON(w, http.StatusNotFound, response.Body{"message": err.Error()})
			return
		}

		h.logger.Error("failed to proxy url", slog.String("error", err.Error()))
		response.JSON(w, http.StatusInternalServerError, response.Body{"message": "failed to proxy url"})

	}

	http.Redirect(w, r, link, http.StatusMovedPermanently)
}

func getUrlFromPayload(r *http.Request) (string, error) {
	var input request.URLInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return "", err
	}

	if err := input.Validate(); err != nil {
		return "", err
	}

	return input.URL, nil
}

func buildShortURL(encoder ServiceEncoder, host string, link *domain.Link) (string, string) {
	code := encoder.Encode(link.ID)
	return code, fmt.Sprintf("http://%s/%s", host, code)
}
