package service

import (
	"errors"
	"github.com/g-stro/content-service/internal/dto"
	"github.com/g-stro/content-service/internal/model"
	"github.com/g-stro/content-service/internal/repository"
	"log/slog"
	"time"
)

type clock func() time.Time

type Service struct {
	repo  repository.ContentRepository
	clock clock
}

func NewContentService(repo repository.ContentRepository, clock clock) *Service {
	if clock == nil {
		clock = time.Now // Default
	}

	return &Service{
		repo:  repo,
		clock: clock,
	}
}

func (s *Service) GetContent() ([]*dto.Content, error) {
	content, err := s.repo.GetAllContent()
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return []*dto.Content{}, nil
	}

	res := make([]*dto.Content, 0)
	for _, c := range content {
		contentDTO, err := s.convertContentModelToDTO(c)
		if err != nil {
			return nil, err
		}
		res = append(res, contentDTO)
	}

	return res, nil
}

func (s *Service) CreateContent(req dto.Content) (*dto.Content, error) {
	content, err := s.convertContentDTOToModel(&req)
	if err != nil {
		return nil, errors.New("failed to convert CreateRequestDTO to model")
	}

	content, err = s.repo.CreateContentWithDetails(content)
	if err != nil {
		return nil, err
	}

	resp, err := s.convertContentModelToDTO(content)
	if err != nil {
		return nil, errors.New("failed to convert model to response DTO")
	}

	return resp, nil
}

// convertContentTypeNameToID converts a content type name string to content type ID integer
func (s *Service) convertContentTypeNameToID(name string) (int, error) {
	ct, err := s.repo.GetContentTypeByName(name)
	if err != nil {
		slog.Error("failed to fetch ContentTypeID", "error", err)
		return 0, err
	}
	return ct.ID, nil
}

// convertContentTypeIDToName converts a content type ID integer to content type string
func (s *Service) convertContentTypeIDToName(id int) (string, error) {
	ct, err := s.repo.GetContentTypeByID(id)
	if err != nil {
		slog.Error("failed to fetch ContentTypeName", "error", err)
		return "", err
	}
	return ct.Name, nil
}

func (s *Service) convertContentDTOToModel(content *dto.Content) (*model.Content, error) {
	if content == nil {
		err := errors.New("request DTO is nil")
		slog.Error("request DTO is nil", "error", err)
		return nil, err
	}

	currTime := s.clock()
	res := &model.Content{
		Title:            content.Title,
		Description:      content.Description,
		CreationDate:     currTime,
		LastModifiedDate: currTime,
	}

	// Convert the res details
	if content.Details != nil {
		for _, d := range content.Details {
			contentTypeID, err := s.convertContentTypeNameToID(d.ContentType)
			if err != nil {
				slog.Error("failed to convert res type to ID", "error", err)
				return nil, err
			}
			detail := model.Details{
				ContentTypeID: contentTypeID,
				Value:         d.Value,
			}
			res.Details = append(res.Details, &detail)
		}
	}

	return res, nil
}

func (s *Service) convertContentModelToDTO(content *model.Content) (*dto.Content, error) {
	if content == nil {
		err := errors.New("content is nil")
		slog.Error("content is nil", "error", err)
		return nil, err
	}

	res := &dto.Content{
		ID:           content.ID,
		Title:        content.Title,
		CreationDate: content.CreationDate,
		Description:  content.Description,
	}

	// Convert the content details
	if content.Details != nil {
		for _, d := range content.Details {
			contentType, err := s.convertContentTypeIDToName(d.ContentTypeID)
			if err != nil {
				slog.Error("failed to convert ID to content type", "error", err)
				return nil, err
			}
			detail := dto.Details{
				ContentType: contentType,
				Value:       d.Value,
			}
			res.Details = append(res.Details, detail)
		}
	}

	return res, nil
}
