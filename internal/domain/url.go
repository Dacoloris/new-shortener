package domain

type URL struct {
	UserID   string `json:"-"`
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

type BatchPostRequest struct {
	CorrelationID string `json:"correlation_id"`
	Original      string `json:"original_url"`
}

type BatchPostResponse struct {
	CorrelationID string `json:"correlation_id"`
	Short         string `json:"short_url"`
}
