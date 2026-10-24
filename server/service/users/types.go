package users

type ImpersonatedByUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserPayload struct {
	ID             int                 `json:"id"`
	Username       string              `json:"username"`
	Email          string              `json:"email"`
	ImpersonatedBy *ImpersonatedByUser `json:"impersonatedBy,omitempty"`
}
