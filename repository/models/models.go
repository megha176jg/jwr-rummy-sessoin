package models

type AuthToken struct {
	AuthToken string `json:"auth_token"`
	Error     error  `json:"error"`
}
