package domain

type URL struct {
	UserID   string `json:"-"`
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}
