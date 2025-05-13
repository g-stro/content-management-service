package service

import (
	"errors"
	"github.com/g-stro/content-service/internal/dto"
	"github.com/g-stro/content-service/internal/model"
	"reflect"
	"testing"
	"time"
)

var fixedTime = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
var testClock = func() time.Time {
	return fixedTime
}

type MockRepository struct {
	MockedContent  []*model.Content
	MockedError    error
	CreatedContent *model.Content

	ContentTypeNameToIDMap map[string]*model.ContentType
	ContentTypeIDToNameMap map[int]*model.ContentType
}

func (m *MockRepository) GetAllContent() ([]*model.Content, error) {
	if m.MockedError != nil {
		return nil, m.MockedError
	}
	return m.MockedContent, nil
}

func (m *MockRepository) CreateContentWithDetails(content *model.Content) (*model.Content, error) {
	if m.MockedError != nil {
		return nil, m.MockedError
	}
	m.CreatedContent = content
	content.ID = 1
	return content, nil
}

func (m *MockRepository) GetContentTypeByName(name string) (*model.ContentType, error) {
	if m.MockedError != nil {
		return nil, m.MockedError
	}
	val, ok := m.ContentTypeNameToIDMap[name]
	if !ok {
		return nil, errors.New("content type not found")
	}
	return val, nil
}

func (m *MockRepository) GetContentTypeByID(id int) (*model.ContentType, error) {
	if m.MockedError != nil {
		return nil, m.MockedError
	}
	val, ok := m.ContentTypeIDToNameMap[id]
	if !ok {
		return nil, errors.New("content type not found")
	}
	return val, nil
}

func TestService_GetContent(t *testing.T) {
	tests := []struct {
		name      string
		repoMock  *MockRepository
		expected  []*dto.Content
		expectErr bool
	}{
		{
			name: "successful fetch",
			repoMock: &MockRepository{
				MockedContent: []*model.Content{
					{
						ID:          1,
						Title:       "Test Title",
						Description: "Test Description",
						Details:     nil,
					},
				},
			},
			expected: []*dto.Content{
				{
					ID:          1,
					Title:       "Test Title",
					Description: "Test Description",
					Details:     nil,
				},
			},
			expectErr: false,
		},
		{
			name: "no content available",
			repoMock: &MockRepository{
				MockedContent: []*model.Content{},
			},
			expected:  []*dto.Content{},
			expectErr: false,
		},
		{
			name: "error while fetching content",
			repoMock: &MockRepository{
				MockedError: errors.New("repository error"),
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewContentService(tt.repoMock, testClock)

			result, err := service.GetContent()

			if (err != nil) != tt.expectErr {
				t.Errorf("GetContent() error = %v, expectErr = %v", err, tt.expectErr)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetContent() got = %v, expected = %v", result, tt.expected)
			}
		})
	}
}

func TestService_CreateContent(t *testing.T) {
	tests := []struct {
		name      string
		input     dto.Content
		repoMock  *MockRepository
		expected  *dto.Content
		expectErr bool
	}{
		{
			name: "successful creation",
			input: dto.Content{
				Title:       "Test Title",
				Description: "Test Description",
			},
			repoMock: &MockRepository{},
			expected: &dto.Content{
				ID:           1,
				Title:        "Test Title",
				Description:  "Test Description",
				CreationDate: fixedTime,
			},
			expectErr: false,
		},
		{
			name: "repository error creating content",
			input: dto.Content{
				Title:       "Test Title",
				Description: "Test Description",
			},
			repoMock: &MockRepository{
				MockedError: errors.New("repository error"),
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewContentService(tt.repoMock, testClock)

			result, err := service.CreateContent(tt.input)

			if (err != nil) != tt.expectErr {
				t.Errorf("CreateContent() error = %v, expectErr = %v", err, tt.expectErr)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CreateContent() got = %v, expected = %v", result, tt.expected)
			}
		})
	}
}
