package model

type UserStruct struct {
	Id            string `json:"id"`
	PhoneNumber   string `json:"phoneNumber"`
	Password      string `json:"password"`
	LoginRetry    int    `json:"login_retry"`
	NextLoginDate string `json:"next_login_date"`
	LastLogin     string `json:"last_login"`
	SessionId     string `json:"session_id"`
	Status        int    `json:"status"`
	CreatedDate   string `json:"created_date"`
}
