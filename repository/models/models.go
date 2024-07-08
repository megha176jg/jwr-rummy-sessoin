package models

type AuthToken struct {
	AuthToken string `json:"auth_token"`
	Error     error  `json:"error"`
}

type DeleteAuthToken struct {
	Error error `json:"error"`
}

type CreateAuthToken struct {
	UserId    string `json:"user_id"`
	AuthToken string `json:"auth_token"`
	Error     string `json:"error"`
}
