package auth

import "time"

//step1 (register): request user dalam bentuk json
type RegisterRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequestPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogoutRequestPayload struct {
	Id          string
	AccessToken string
	Exp         time.Time
}

type RequestPayLoadSuperAdmin struct {
	Email    string
	Password string
}
