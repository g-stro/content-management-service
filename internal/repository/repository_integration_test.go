//go:build integration

package repository

import (
	"fmt"
	"github.com/g-stro/content-service/database"
	"github.com/g-stro/content-service/internal/model"
	"reflect"
	"strings"
	"testing"
	"time"
)

var (
	staticTimestamp = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	testName        = "test name"
	testDescription = "test description"
)

func TestPostgresContentRepository_GetAllContent(t *testing.T) {
	conn, err := database.NewConnection()
	if err != nil {
		t.Fatalf("failed to establish database connection: %v", err)
	}
	defer conn.DB.Close()

	repo := NewPostgresContentRepository(conn)

	tests := []struct {
		name     string
		setup    func() error
		expected []*model.Content
		wantErr  bool
	}{
		{
			name: "successful fetch",
			setup: func() error {
				var id int
				err := conn.DB.QueryRow(
					`INSERT INTO content (name, description, creation_date, last_modified_date) 
					VALUES ($1, $2, $3, $4) 
					RETURNING id`,
					testName, testDescription, staticTimestamp, staticTimestamp).Scan(&id)
				if err != nil {
					return err
				}

				var detailsID int
				err = conn.DB.QueryRow(
					`INSERT INTO content_details (content_id, content_type_id, value)
					VALUES ($1, $2, $3)
					RETURNING id`,
					id, 1, "test text").Scan(&detailsID)
				if err != nil {
					return err
				}

				return err
			},
			expected: []*model.Content{
				{ID: 1, Name: testName, Description: testDescription, CreationDate: staticTimestamp,
					LastModifiedDate: staticTimestamp, Details: []*model.Details{{ID: 1, ContentID: 1, ContentTypeID: 1, Value: "test text"}}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.setup(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			// Clean the database
			defer func() {
				if _, err := conn.DB.Exec("DELETE FROM content_details; DELETE FROM content;"); err != nil {
					t.Fatalf("Failed to clean up database: %v", err)
				}
			}()

			content, err := repo.GetAllContent()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllContent() error = %v, expected error = %v", err, tt.wantErr)
				return
			}

			if !isContentSliceEqual(content, tt.expected) {
				t.Errorf("GetAllContent() got: \n%+v\nexpected:\n%+v", printSlice(content), printSlice(tt.expected))
			}
		})
	}
}

func TestPostgresContentRepository_CreateContentWithDetails(t *testing.T) {
	conn, err := database.NewConnection()
	if err != nil {
		t.Fatalf("failed to establish database connection: %v", err)
	}
	defer conn.DB.Close()

	repo := NewPostgresContentRepository(conn)

	tests := []struct {
		name     string
		input    *model.Content
		expected *model.Content
		wantErr  bool
	}{
		{
			name: "successful creation",
			input: &model.Content{Name: testName, Description: testDescription, CreationDate: staticTimestamp,
				LastModifiedDate: staticTimestamp, Details: []*model.Details{{ContentTypeID: 1, Value: "test text"}}},
			expected: &model.Content{ID: 2, Name: testName, Description: testDescription, CreationDate: staticTimestamp,
				LastModifiedDate: staticTimestamp, Details: []*model.Details{{ID: 2, ContentID: 2, ContentTypeID: 1, Value: "test text"}}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call repository method
			content, err := repo.CreateContentWithDetails(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateContent() error = %v, expected error = %v", err, tt.wantErr)
				return
			}

			// Clean the database
			defer func() {
				if _, err := conn.DB.Exec("DELETE FROM content_details; DELETE FROM content;"); err != nil {
					t.Fatalf("Failed to clean up database: %v", err)
				}
			}()

			// Compare Details slice
			if !isSliceEqual(tt.expected.Details, content.Details, func(dExp, dAct *model.Details) bool {
				if dExp == nil || dAct == nil {
					return dExp == dAct
				}
				return reflect.DeepEqual(*dExp, *dAct)
			}) {
				t.Errorf("CreateContent() got: \n%+v\nexpected:\n%+v", fmt.Sprintf("%+v\n", *content), fmt.Sprintf("%+v\n", tt.expected))
			}

			// Compare all other fields
			tt.expected.Details = nil
			content.Details = nil
			if !reflect.DeepEqual(*content, *tt.expected) {
				t.Errorf("CreateContent() got: \n%+v\nexpected:\n%+v", fmt.Sprintf("%+v\n", *content), fmt.Sprintf("%+v\n", tt.expected))
			}
		})
	}
}

// TestPostgresContentRepository_GetContentTypeByName tests the GetContentTypeByName method of PostgresContentRepository
// with data already seeded from sql/sql.sql
func TestPostgresContentRepository_GetContentTypeByName(t *testing.T) {
	conn, err := database.NewConnection()
	if err != nil {
		t.Fatalf("failed to establish database connection: %v", err)
	}
	defer conn.DB.Close()

	repo := NewPostgresContentRepository(conn)

	// Define test cases
	tests := []struct {
		name            string
		contentTypeName string
		expectedID      int
		wantErr         bool
	}{
		{
			name:            "successful fetch",
			contentTypeName: "text",
			expectedID:      1,
			wantErr:         false,
		},
		{
			name:            "content type not found",
			contentTypeName: "Unknown",
			expectedID:      0, // Expect nil response
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call repository method
			contentType, err := repo.GetContentTypeByName(tt.contentTypeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContentTypeByName() error = %v, expected error = %v", err, tt.wantErr)
				return
			}

			// Check nil case
			if contentType == nil && tt.expectedID != 0 {
				t.Errorf("GetContentTypeByName() got = nil, expected = %v", tt.expectedID)
				return
			}

			if contentType != nil && contentType.ID != tt.expectedID {
				t.Errorf("GetContentTypeByName() got = %v, expected = %v", contentType.ID, tt.expectedID)
			}
		})
	}
}

// TestPostgresContentRepository_GetContentTypeByID tests the GetContentTypeByID method of PostgresContentRepository
// with data already seeded from sql/sql.sql
func TestPostgresContentRepository_GetContentTypeByID(t *testing.T) {
	conn, err := database.NewConnection()
	if err != nil {
		t.Fatalf("failed to establish database connection: %v", err)
	}
	defer conn.DB.Close()

	repo := NewPostgresContentRepository(conn)

	// Define test cases
	tests := []struct {
		name         string
		id           int
		expectedName string
		wantErr      bool
	}{
		{
			name:         "successful fetch",
			id:           1,
			expectedName: "text",
			wantErr:      false,
		},
		{
			name:         "content type not found",
			id:           0,
			expectedName: "", // Expect nil response
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call repository method
			contentType, err := repo.GetContentTypeByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContentTypeByID() error = %v, expected error = %v", err, tt.wantErr)
				return
			}

			// Check nil case
			if contentType == nil && tt.expectedName != "" {
				t.Errorf("GetContentTypeByID() got = nil, expected = %v", tt.expectedName)
				return
			}

			if contentType != nil && contentType.Name != tt.expectedName {
				t.Errorf("GetContentTypeByID() got = %v, expected = %v", contentType.Name, tt.expectedName)
			}
		})
	}
}

func printSlice[T any](items []*T) string {
	var builder strings.Builder
	for _, item := range items {
		builder.WriteString(fmt.Sprintf("%+v\n", *item))
	}
	return builder.String()
}

func isSliceEqual[T any](expected, actual []T, compare func(e, a T) bool) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i, exp := range expected {
		if !compare(exp, actual[i]) {
			return false
		}
	}

	return true
}

func isContentSliceEqual(expected, actual []*model.Content) bool {
	return isSliceEqual(expected, actual, func(exp, act *model.Content) bool {
		if exp == nil || act == nil {
			return exp == act
		}

		// Compare Details slice
		if !isSliceEqual(exp.Details, act.Details, func(dExp, dAct *model.Details) bool {
			if dExp == nil || dAct == nil {
				return dExp == dAct
			}
			return reflect.DeepEqual(*dExp, *dAct)
		}) {
			return false
		}

		// Compare all other fields
		exp.Details = nil
		act.Details = nil
		return reflect.DeepEqual(*exp, *act)
	})
}
