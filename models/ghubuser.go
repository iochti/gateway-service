package models

type GhubUser struct {
	Login     string `json:"login"`
	Id        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}
