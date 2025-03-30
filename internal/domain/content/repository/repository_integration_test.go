//go:build integration

package repository

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/g-stro/content-service/internal/database"
	"github.com/g-stro/content-service/internal/domain/content/model"
)

var (
	staticTimestamp = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	testTitle       = "test title"
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
					`INSERT INTO content (title, description, creation_date, last_modified_date) 
					VALUES ($1, $2, $3, $4) 
					RETURNING id`,
					testTitle, testDescription, staticTimestamp, staticTimestamp).Scan(&id)
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
				{ID: 1, Title: testTitle, Description: testDescription, CreationDate: staticTimestamp,
					LastModifiedDate: staticTimestamp, Details: []*model.Detail{{ID: 1, ContentID: 1, ContentTypeID: 1, Value: "test text"}}},
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
			input: &model.Content{Title: testTitle, Description: testDescription, CreationDate: staticTimestamp,
				LastModifiedDate: staticTimestamp, Details: []*model.Detail{{ContentTypeID: 1, Value: "test text"}}},
			expected: &model.Content{ID: 2, Title: testTitle, Description: testDescription, CreationDate: staticTimestamp,
				LastModifiedDate: staticTimestamp, Details: []*model.Detail{{ID: 2, ContentID: 2, ContentTypeID: 1, Value: "test text"}}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call repository method
			content, err := repo.CreateContentWithDetails(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateContentWithDetails() error = %v, expected error = %v", err, tt.wantErr)
				return
			}

			// Clean the database
			defer func() {
				if _, err := conn.DB.Exec("DELETE FROM content_details; DELETE FROM content;"); err != nil {
					t.Fatalf("Failed to clean up database: %v", err)
				}
			}()

			// Compare Details slice
			if !isSliceEqual(tt.expected.Details, content.Details, func(dExp, dAct *model.Detail) bool {
				if dExp == nil || dAct == nil {
					return dExp == dAct
				}
				return reflect.DeepEqual(*dExp, *dAct)
			}) {
				t.Errorf("CreateContentWithDetails() got: \n%+v\nexpected:\n%+v", fmt.Sprintf("%+v\n", *content), fmt.Sprintf("%+v\n", tt.expected))
			}

			// Compare all other fields
			tt.expected.Details = nil
			content.Details = nil
			if !reflect.DeepEqual(*content, *tt.expected) {
				t.Errorf("CreateContentWithDetails() got: \n%+v\nexpected:\n%+v", fmt.Sprintf("%+v\n", *content), fmt.Sprintf("%+v\n", tt.expected))
			}
		})
	}
}

// TestPostgresContentRepository_GetContentTypeID tests the GetContentTypeID method of PostgresContentRepository
// with data already seeded from sql/sql.sql
func TestPostgresContentRepository_GetContentTypeID(t *testing.T) {
	conn, err := database.NewConnection()
	if err != nil {
		t.Fatalf("failed to establish database connection: %v", err)
	}
	defer conn.DB.Close()

	repo := NewPostgresContentRepository(conn)

	// Define test cases
	tests := []struct {
		name        string
		contentType string
		expectedID  int
		wantErr     bool
	}{
		{
			name:        "successful fetch",
			contentType: "text",
			expectedID:  1,
			wantErr:     false,
		},
		{
			name:        "content type not found",
			contentType: "Unknown",
			expectedID:  0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call repository method
			id, err := repo.GetContentTypeID(tt.contentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContentTypeID() error = %v, expected error = %v", err, tt.wantErr)
				return
			}

			if id != tt.expectedID {
				t.Errorf("GetContentTypeID() got = %v, expected = %v", id, tt.expectedID)
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
		if !isSliceEqual(exp.Details, act.Details, func(dExp, dAct *model.Detail) bool {
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
