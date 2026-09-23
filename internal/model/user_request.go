package model

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

func (r CreateUserRequest) Validate() string {
	if r.Username == "" {
		return "Username is required"
	}

	if r.Password == "" {
		return "Password is required"
	}

	if len(r.Password) < 8 {
		return "Password must be at least 8 characters"
	}

	if r.Role == "" {
		return "Role is required"
	}

	if r.Role != "viewer" &&
		r.Role != "staff" &&
		r.Role != "admin" {
		return "Invalid role"
	}

	if r.Email == "" {
		return "Email is required"
	}

	return ""
}
