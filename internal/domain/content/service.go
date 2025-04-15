package content

import (
	"encoding/json"
	"errors"
	"github.com/g-stro/content-service/internal/domain/content/dto"
	"github.com/g-stro/content-service/internal/domain/content/model"
	"github.com/g-stro/content-service/internal/domain/content/repository"
	"github.com/g-stro/content-service/internal/response"
	"log/slog"
	"net/http"
	"time"
)

type Service struct {
	ctx  *http.ServeMux
	repo *repository.PostgresContentRepository
}

func (s *Service) contentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.getContent(w, r)
	case "POST":
		s.createContentWithDetails(w, r)
	//case "PUT":
	//s.updateContent(w, r)
	//case "DELETE":
	//s.deleteContent(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func NewContentService(c *http.ServeMux, repo *repository.PostgresContentRepository) *Service {
	s := &Service{ctx: c, repo: repo}
	s.ctx.HandleFunc("/content", s.contentHandler)
	return s
}

func (s *Service) getContent(w http.ResponseWriter, r *http.Request) {
	content, err := s.repo.GetAllContent()
	if err != nil {
		response.HttpError(w, err, http.StatusInternalServerError, "failed to retrieve content")
		return
	}

	if len(content) == 0 {
		response.HttpSuccess(w, map[string]interface{}{
			"content": []dto.ContentWithDetailsResponse{},
		}, http.StatusOK, "No content available")
		return
	}

	resp := make([]*dto.ContentWithDetailsResponse, 0)
	for _, c := range content {
		contentDTO, err := s.convertContentModelToResponse(c)
		if err != nil {
			response.HttpError(w, err, http.StatusInternalServerError, "failed to convert content to response")
			return
		}
		resp = append(resp, contentDTO)
	}

	result := struct {
		Content []*dto.ContentWithDetailsResponse `json:"content"`
	}{
		Content: resp,
	}

	response.HttpSuccess(w, result, http.StatusOK, "content retrieved successfully")
}

func (s *Service) createContentWithDetails(w http.ResponseWriter, r *http.Request) {
	// Validate the request
	var req dto.ContentWithDetailsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.HttpFail(
			w, "invalid request data format", http.StatusBadRequest, "invalid request data format")
		return
	}

	// Convert request to domain model
	content, err := s.convertCreateContentRequestToModel(&req)
	if err != nil {
		response.HttpFail(w, "failed to convert CreateRequestDTO to model", http.StatusBadRequest, "failed to convert CreateRequestDTO to model")
		return
	}

	// Save the new content
	content, err = s.repo.CreateContentWithDetails(content)
	if err != nil {
		response.HttpError(w, err, http.StatusInternalServerError, "failed to create content")
		return
	}

	// Convert the saved content to response
	resp, err := s.convertContentModelToResponse(content)
	if err != nil {
		response.HttpError(w, err, http.StatusInternalServerError, "failed to build response")
		return
	}

	response.HttpSuccess(w, resp, http.StatusOK, "content created successfully")
}

// convertContentTypeToID converts a content type string to content type ID integer
func (s *Service) convertContentTypeToID(contentType string) (int, error) {
	contentTypeID, err := s.repo.GetContentTypeID(contentType)
	if err != nil {
		slog.Error("failed to fetch ContentTypeID", "error", err)
		return 0, err
	}
	return contentTypeID, nil
}

// convertContentTypeIDToName converts a content type ID integer to content type string
func (s *Service) convertContentTypeIDToName(contentTypeID int) (string, error) {
	contentTypeName, err := s.repo.GetContentTypeName(contentTypeID)
	if err != nil {
		slog.Error("failed to fetch ContentTypeName", "error", err)
		return "", err
	}
	return contentTypeName, nil
}

func (s *Service) convertCreateContentRequestToModel(req *dto.ContentWithDetailsRequest) (*model.Content, error) {
	if req == nil {
		err := errors.New("request DTO is nil")
		slog.Error("request DTO is nil", "error", err)
		return nil, err
	}

	currTime := time.Now()
	content := &model.Content{
		Title:            req.Title,
		Description:      req.Description,
		CreationDate:     currTime,
		LastModifiedDate: currTime,
	}

	// Convert the content details
	if req.Details != nil {
		for _, d := range req.Details {
			contentTypeID, err := s.convertContentTypeToID(d.ContentType)
			if err != nil {
				slog.Error("failed to convert content type to ID", "error", err)
				return nil, err
			}
			detail := model.Detail{
				ContentTypeID: contentTypeID,
				Value:         d.Value,
			}
			content.Details = append(content.Details, &detail)
		}
	}

	return content, nil
}

func (s *Service) convertContentModelToResponse(content *model.Content) (*dto.ContentWithDetailsResponse, error) {
	if content == nil {
		err := errors.New("content is nil")
		slog.Error("content is nil", "error", err)
		return nil, err
	}

	resp := &dto.ContentWithDetailsResponse{
		Title:       content.Title,
		Description: content.Description,
	}

	// Convert the content details
	if content.Details != nil {
		for _, d := range content.Details {
			contentType, err := s.convertContentTypeIDToName(d.ContentTypeID)
			if err != nil {
				slog.Error("failed to convert ID to content type", "error", err)
				return nil, err
			}
			detail := dto.Detail{
				ContentType: contentType,
				Value:       d.Value,
			}
			resp.Details = append(resp.Details, detail)
		}
	}

	return resp, nil
}
