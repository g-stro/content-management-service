package dto

import "time"

type Content struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreationDate time.Time `json:"creation_date"`
	Details      []Details `json:"details"`
}

type Details struct {
	ContentType string `json:"content_type"`
	Value       string `json:"value"`
}
