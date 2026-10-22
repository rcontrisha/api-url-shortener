package dto

import "time"

type ShortenLink struct {
	OriginalUrl string `json:"original_url" binding:"required,url"`
	Slug        string `json:"slug,omitempty" binding:"omitempty,max=50"`
}

type LinkItemResponse struct {
	ID          int       `json:"id"`
	Slug        string    `json:"slug"`
	OriginalUrl string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaginationMeta struct {
	TotalActive int `json:"total_active"`
	Page        int `json:"page"`
	Limit       int `json:"limit"`
	TotalPages  int `json:"total_pages"`
}

type GetUserLinksResponse struct {
	Data []LinkItemResponse `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}
