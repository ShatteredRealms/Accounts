package model

// LoginRequest The model representing how a user would sign-in to the application via a form or JSON data. Both Email
// and Password are requried.
type LoginRequest struct {
	// Email The username supplied for login
	Email string `form:"email" json:"email" binding:"required"`

	// Password the password supplied for login
	Password string `form:"password" json:"password" binding:"required"`
}

// LoginResponse The data to respond with upon a successful login attempt
type LoginResponse struct {
	// Token the JWT token for the authenticated user
	Token string `json:"token"`
}
