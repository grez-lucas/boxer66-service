package users

type APIResponseStatus string

const (
	APIResponseStatusSuccess = "success"
	APIResponseStatusError   = "error"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserID int32  `json:"user_id"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmailRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type VerifyEmailResponse struct {
	Token  string `json:"token"`
	UserID int32  `json:"user_id"`
}

type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
