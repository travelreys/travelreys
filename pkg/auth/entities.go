package auth

type User struct {
	ID     string            `json:"id"`
	Email  string            `json:"email"`
	Labels map[string]string `json:"labels"`
}
