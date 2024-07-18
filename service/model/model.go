package model

// type Create struct {
// 	AuthToken string `json:"auth_token"`
// 	Err       error  `json:"error"`
// }

// type Validate struct {
// 	IsTokenValid bool  `json:"is_token_valid"`
// 	Err          error `json:"error"`
// }

//	type Invalidate struct {
//		Err error `json:"error"`
//	}
type TokenCreate struct {
	Error Err `json:"error"`
}
type TokenGet struct {
	AuthToken string `json:"auth_token"`
	Error     Err    `json:"error"`
}
type TokenValidation struct {
	IsTokenValid bool `json:"is_token_valid"`
	Error        Err  `json:"error"`
}
type TokenDeletion struct {
	Error Err `json:"error"`
}
type Err struct {
	Message string `json:"error_message"`
}
