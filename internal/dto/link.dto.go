package dto

type ShortenLink struct {
	OriginalUrl string `json:"original_url" binding:"required,url"`
	Slug        string `json:"slug,omitempty" binding:"omitempty,max=50"`
}