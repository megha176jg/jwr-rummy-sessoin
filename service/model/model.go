package model

type Create struct {
	AuthToken string `json:"auth_token"`
	Err       error  `json:"error"`
}

type Validate struct {
	IsTokenValid bool  `json:"is_token_valid"`
	Err          error `json:"error"`
}

type Invalidate struct {
	Err error `json:"error"`
}
