package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/response"
	"BE_GKMI_NTC_GO/internal/service"
	"BE_GKMI_NTC_GO/internal/token"
)

// Login godoc
// @Summary Login
// @Description Login menggunakan username dan password untuk mendapatkan access token dan refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/login [post]
func Login(authService AuthServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request model.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, refreshToken, err := authService.Login(r.Context(), request)
		if err != nil {
			if errors.Is(err, service.ErrInvalidCredentials) {
				response.Error(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}
			response.Error(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		accessToken, err := token.GenerateAccessToken(user.ID, user.Username, user.Role)
		if err != nil {
			http.Error(w, "failed to generate access token", http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, map[string]string{
			"access_token":  accessToken,
			"token_type":    "Bearer",
			"refresh_token": refreshToken,
		})
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Generate a new access token using a refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/refresh [post]
func Refresh(authService AuthServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request model.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, refreshToken, err := authService.Refresh(r.Context(), request.RefreshToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		accessToken, err := token.GenerateAccessToken(user.ID, user.Username, user.Role)
		if err != nil {
			http.Error(w, "failed to generate access token", http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, map[string]string{
			"access_token":  accessToken,
			"token_type":    "Bearer",
			"refresh_token": refreshToken,
		})
	})
}

// Logout godoc
// @Summary Logout
// @Description Revoke refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/auth/logout [post]
func Logout(authService AuthServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request model.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := authService.Logout(r.Context(), request.RefreshToken); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		response.JSON(w, http.StatusOK, map[string]string{
			"message": "logout successful",
		})
	})
}

func Me(userID int, username string, role string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"user_id":  userID,
			"username": username,
			"role":     role,
		})
	})
}
