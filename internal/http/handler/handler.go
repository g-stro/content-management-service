package handler

import (
	"encoding/json"
	"github.com/g-stro/content-service/internal/dto"
	"github.com/g-stro/content-service/internal/http/response"
	"github.com/g-stro/content-service/internal/service"
	"net/http"
)

type Handler struct {
	svc *service.Service
}

func NewContentHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/content", h.handleContentRequests)
}

func (h *Handler) handleContentRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getContent(w, r)
	case http.MethodPost:
		h.createContent(w, r)
	default:
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getContent(w http.ResponseWriter, r *http.Request) {
	content, err := h.svc.GetContent()
	if err != nil {
		response.HttpError(w, err, http.StatusInternalServerError, "failed to retrieve content")
		return
	}

	if len(content) == 0 {
		response.HttpSuccess(w, map[string]interface{}{
			"content": []response.GetContent{},
		}, http.StatusOK, "No content available")
		return
	}

	contentResp := make([]response.GetContent, 0)
	for _, c := range content {
		details := make([]response.Details, 0)
		for _, d := range c.Details {
			dr := response.Details{
				ContentType: d.ContentType,
				Value:       d.Value,
			}
			details = append(details, dr)
		}

		cr := response.GetContent{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			Details:     details,
		}
		contentResp = append(contentResp, cr)
	}

	resp := struct {
		Content []response.GetContent `json:"content"`
	}{
		Content: contentResp,
	}

	response.HttpSuccess(w, resp, http.StatusOK, "content retrieved successfully")
}

func (h *Handler) createContent(w http.ResponseWriter, r *http.Request) {
	var req dto.Content
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.HttpFail(
			w, "invalid request body", http.StatusBadRequest, "invalid request body")
		return
	}

	content, err := h.svc.CreateContent(req)
	if err != nil {
		response.HttpError(w, err, http.StatusInternalServerError, "failed to create content")
		return
	}

	formattedCreationDate := content.CreationDate.Format("2006-01-02 15:04:05")

	resp := response.CreateContent{
		ID:           content.ID,
		Name:         content.Name,
		CreationDate: formattedCreationDate,
	}

	response.HttpSuccess(w, resp, http.StatusCreated, "content created successfully")
}
