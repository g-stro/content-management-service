package repository

import (
	"database/sql"
	"github.com/g-stro/content-service/internal/database"
	"github.com/g-stro/content-service/internal/domain/content/model"
	"log/slog"
)

type ContentRepository interface {
	GetAllContent() ([]*model.Content, error)
	CreateContentWithDetails(content *model.Content) (*model.Content, error)
	GetContentTypeID(contentType string) (string, error)
}

type PostgresContentRepository struct {
	conn *database.Connection
}

func NewPostgresContentRepository(c *database.Connection) *PostgresContentRepository {
	return &PostgresContentRepository{conn: c}
}

func (r *PostgresContentRepository) GetAllContent() ([]*model.Content, error) {
	query := `SELECT c.id, c.title, c.description, c.creation_date, c.last_modified_date,
                 cd.id, cd.content_id, cd.content_type_id, cd.value
                 FROM content c 
                 JOIN content_details cd ON c.id = cd.content_id
                 JOIN content_type ct ON ct.id = cd.content_type_id`

	rows, err := r.conn.DB.Query(query)
	if err != nil {
		slog.Error("failed to execute query", "error", err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}(rows)

	var contentsMap = make(map[any]*model.Content)
	for rows.Next() {
		var content model.Content
		var contentDetail model.Detail
		err = rows.Scan(
			&content.ID, &content.Title, &content.Description, &content.CreationDate, &content.LastModifiedDate,
			&contentDetail.ID, &contentDetail.ContentID, &contentDetail.ContentTypeID, &contentDetail.Value)
		if err != nil {
			slog.Error("failed to scan rows into content and contentDetail structures", "error", err)
			return nil, err
		}

		// Normalize times to UTC
		content.CreationDate = content.CreationDate.UTC()
		content.LastModifiedDate = content.LastModifiedDate.UTC()

		if _, exists := contentsMap[content.ID]; !exists {
			content.Details = make([]*model.Detail, 0)
			contentsMap[content.ID] = &content
		}
		contentsMap[content.ID].Details = append(contentsMap[content.ID].Details, &contentDetail)
	}

	var result []*model.Content
	for _, content := range contentsMap {
		result = append(result, content)
	}

	return result, nil
}

func (r *PostgresContentRepository) CreateContentWithDetails(content *model.Content) (*model.Content, error) {
	tx, err := r.conn.DB.Begin()
	if err != nil {
		slog.Error("failed to start the transaction", "error", err)
		return nil, err
	}

	defer func() {
		if err != nil {
			slog.Error("transaction error", "error", err)
			err := tx.Rollback()
			if err != nil {
				slog.Error("failed to roll back transaction", "error", err)
			}
		}
	}()

	stmtContent := `
        INSERT INTO content (title, description, creation_date, last_modified_date)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	var id int
	err = tx.QueryRow(
		stmtContent, content.Title, content.Description, content.CreationDate, content.LastModifiedDate).Scan(&id)
	if err != nil {
		slog.Error("failed to execute query and scan result", "error", err)
		return nil, err
	}

	stmtDetails := `
	   INSERT INTO content_details (content_id, content_type_id, value)
	   VALUES ($1, $2, $3)
	   RETURNING id`

	for _, cd := range content.Details {
		cd.ContentID = id
		var detailsID int
		err = tx.QueryRow(stmtDetails, cd.ContentID, cd.ContentTypeID, cd.Value).Scan(&detailsID)
		if err != nil {
			slog.Error("failed to execute details query or scan result", "error", err)
			return nil, err
		}
		cd.ID = detailsID // Set the content details ID after creation.
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		slog.Error("failed to commit the transaction", "error", err)
		return nil, err
	}

	content.ID = id // Set the content ID after creation.

	return content, nil
}

func (r *PostgresContentRepository) GetContentTypeID(contentType string) (int, error) {
	var id int
	stmt := "SELECT id FROM content_type WHERE name = $1"
	err := r.conn.DB.QueryRow(stmt, contentType).Scan(&id)
	if err != nil {
		slog.Error("failed to fetch ContentTypeID", "error", err)
		return 0, err
	}
	return id, nil
}
