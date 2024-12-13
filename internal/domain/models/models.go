package models

import "github.com/kartochnik010/outstaff-task/internal/pkg/lib_time"

type SearchMetadata struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`

	ID          uint64            `json:"id,omitempty"`
	Group       string            `json:"group,omitempty"`
	Song        string            `json:"song,omitempty"`
	Link        string            `json:"link,omitempty"`
	Text        string            `json:"text,omitempty"`
	ReleaseDate *lib_time.IntDate `json:"releaseDate,omitempty"`
}

type Music struct {
	ID          uint64           `json:"id"`
	Group       string           `json:"group"`
	Song        string           `json:"song"`
	Link        string           `json:"link"`
	Text        string           `json:"text"`
	ReleaseDate lib_time.IntDate `json:"releaseDate"`
}

// as: github.com/Kartochnik010/lib_time
// as: github.com/kartochnik010/lib_time
