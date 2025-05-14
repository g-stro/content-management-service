package model

import "time"

type Content struct {
	ID               int       `db:"id"`
	Name             string    `db:"name"`
	Description      string    `db:"description"`
	CreationDate     time.Time `db:"creation_date"`
	LastModifiedDate time.Time `db:"last_modified_date"`
	Details          []*Details
}

type Details struct {
	ID            int    `db:"id"`
	ContentID     int    `db:"content_id"`
	ContentTypeID int    `db:"content_type_id"`
	Value         string `db:"value"`
}

type ContentType struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
