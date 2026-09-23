package model

type JemaatPaginationResponse struct {
	Items      []Jemaat `json:"items"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	Total      int      `json:"total"`
	TotalPages int      `json:"total_pages"`
}
