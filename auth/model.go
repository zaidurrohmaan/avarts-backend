package auth

type GoogleLoginRequest struct {
	IdToken string `json:"id_token"`
}

type LoginResponse struct {
	UserID uint   `json:"user_id"`
	IsNew  bool   `json:"is_new"`
	Token  string `json:"token"`
}
