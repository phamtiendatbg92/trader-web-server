package dtos

type GetUserInfoResponse struct {
	Meta *Meta     `json:"meta"`
	Data *UserInfo `json:"user"`
}

type UserInfo struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"user_name"`
}
