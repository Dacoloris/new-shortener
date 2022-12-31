package domain

type URL struct {
	UserID   string `json:"-"`
	Original string `json:"original"`
	Short    string `json:"short"`
}
