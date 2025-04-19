package auth

type GoogleLoginRequest struct {
	IdToken    string `json:"id_token"`
	Name       string `json:"name"`
	AvatarURL  string `json:"avatar_url"`
}