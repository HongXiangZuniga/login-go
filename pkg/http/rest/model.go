package http

type LoginRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}
