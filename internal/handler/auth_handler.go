package handler

import (
	"encoding/json"
	"net/http"

	"BE_GKMI_NTC_GO/internal/middleware"
	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/response"
	"BE_GKMI_NTC_GO/internal/service"
	"BE_GKMI_NTC_GO/internal/token"
)

func Login(authService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var request model.LoginRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		user, err := authService.Login(r.Context(), request)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		accessToken, err := token.GenerateAccessToken(
			user.ID,
			user.Username,
			user.Role,
		)
		if err != nil {
			response.Error(
				w,
				http.StatusInternalServerError,
				"Failed to generate access token",
			)
			return
		}

		_ = response.JSON(w, http.StatusOK, map[string]string{
			"access_token": accessToken,
			"token_type":   "Bearer",
		})
	}
}

func Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "User ID not found")
		return
	}

	username, ok := middleware.GetUsername(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Username not found")
		return
	}

	role, ok := middleware.GetRole(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Role not found")
		return
	}

	_ = response.JSON(w, http.StatusOK, map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"role":     role,
	})
}
