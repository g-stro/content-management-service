package service

import (
	"errors"
	"github.com/g-stro/content-service/internal/domain/content/dto"
	"github.com/g-stro/content-service/internal/domain/content/model"
	"github.com/g-stro/content-service/internal/domain/content/repository"
	"log/slog"
	"time"
)

type Service struct {
	repo *repository.PostgresContentRepository
}

func NewContentService(repo *repository.PostgresContentRepository) *Service {
	s := &Service{repo: repo}
	return s
}

func (s *Service) GetContent() ([]*dto.Content, error) {
	content, err := s.repo.GetAllContent()
	if err != nil {
		//response.HttpError(w, err, http.StatusInternalServerError, "failed to retrieve content")
		return nil, err
	}

	if len(content) == 0 {
		return []*dto.Content{}, nil
	}

	res := make([]*dto.Content, 0)
	for _, c := range content {
		contentDTO, err := s.convertContentModelToResponse(c)
		if err != nil {
			//response.HttpError(w, err, http.StatusInternalServerError, "failed to convert content to response")
			return nil, err
		}
		res = append(res, contentDTO)
	}

	return res, nil
}

func (s *Service) CreateContentWithDetails(req dto.Content) (*dto.Content, error) {
	content, err := s.convertCreateContentRequestToModel(&req)
	if err != nil {
		return nil, errors.New("failed to convert CreateRequestDTO to model")
	}

	content, err = s.repo.CreateContentWithDetails(content)
	if err != nil {
		return nil, err
	}

	resp, err := s.convertContentModelToResponse(content)
	if err != nil {
		return nil, errors.New("failed to convert model to response DTO")
	}

	return resp, nil
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

func (s *Service) convertCreateContentRequestToModel(req *dto.Content) (*model.Content, error) {
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

func (s *Service) convertContentModelToResponse(content *model.Content) (*dto.Content, error) {
	if content == nil {
		err := errors.New("content is nil")
		slog.Error("content is nil", "error", err)
		return nil, err
	}

	resp := &dto.Content{
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
