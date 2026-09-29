package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/repository"
	"BE_GKMI_NTC_GO/internal/response"
	"BE_GKMI_NTC_GO/internal/security"

	"github.com/jackc/pgx/v5/pgconn"
)

// Users godoc
// @Summary List and create users
// @Description Get all users or create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.User
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users [get]
// @Router /api/users [post]
func Users(userRepository *repository.UserRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:

			users, err := userRepository.GetUsers(r.Context())
			if err != nil {
				response.Error(
					w,
					http.StatusInternalServerError,
					"Failed to get users",
				)
				return
			}

			_ = response.JSON(
				w,
				http.StatusOK,
				users,
			)

		case http.MethodPost:

			var request model.CreateUserRequest

			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				response.Error(
					w,
					http.StatusBadRequest,
					"Invalid JSON body",
				)
				return
			}

			validationError := request.Validate()
			if validationError != "" {
				response.Error(
					w,
					http.StatusBadRequest,
					validationError,
				)
				return
			}

			passwordHash, err := security.HashPassword(request.Password)
			if err != nil {
				response.Error(
					w,
					http.StatusInternalServerError,
					"Failed to hash password",
				)
				return
			}

			user := model.User{
				Username:     request.Username,
				PasswordHash: passwordHash,
				Role:         request.Role,
				IsActive:     true,
				Email:        request.Email,
			}

			err = userRepository.CreateUser(
				r.Context(),
				&user,
			)
			if err != nil {
				log.Printf("CREATE USER ERROR: %v", err)

				var pgErr *pgconn.PgError

				if errors.As(err, &pgErr) {
					if pgErr.Code == "23505" {
						response.Error(
							w,
							http.StatusConflict,
							"Username or email already exists",
						)
						return
					}
				}

				response.Error(
					w,
					http.StatusInternalServerError,
					"Failed to create user",
				)
				return
			}

			user.PasswordHash = ""

			_ = response.JSON(
				w,
				http.StatusCreated,
				user,
			)

		default:

			response.Error(
				w,
				http.StatusMethodNotAllowed,
				"Method not allowed",
			)
		}
	}
}

// GetUserByID godoc
// @Summary Get, update, or delete a user
// @Description Get, update, or delete a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users/{id} [get]
// @Router /api/users/{id} [put]
// @Router /api/users/{id} [delete]
func GetUserByID(userRepository *repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idText := strings.TrimPrefix(
			r.URL.Path,
			"/api/users/",
		)

		userID, err := strconv.Atoi(idText)
		if err != nil {
			response.Error(
				w,
				http.StatusBadRequest,
				"Invalid user ID",
			)
			return
		}

		switch r.Method {
		case http.MethodGet:
			user, err := userRepository.GetUserByID(
				r.Context(),
				userID,
			)
			if err != nil {
				response.Error(
					w,
					http.StatusNotFound,
					"User not found",
				)
				return
			}

			_ = response.JSON(
				w,
				http.StatusOK,
				user,
			)

		case http.MethodPut:
			var request struct {
				Username string `json:"username"`
				Role     string `json:"role"`
				IsActive bool   `json:"is_active"`
				Email    string `json:"email"`
			}

			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				response.Error(
					w,
					http.StatusBadRequest,
					"Invalid request body",
				)
				return
			}

			user := &model.User{
				ID:       userID,
				Username: request.Username,
				Role:     request.Role,
				IsActive: request.IsActive,
				Email:    request.Email,
			}

			err = userRepository.UpdateUser(
				r.Context(),
				user,
			)
			if err != nil {
				response.Error(
					w,
					http.StatusInternalServerError,
					"Failed to update user",
				)
				return
			}

			_ = response.JSON(
				w,
				http.StatusOK,
				user,
			)

		case http.MethodDelete:
			err := userRepository.DeleteUser(
				r.Context(),
				userID,
			)
			if err != nil {
				response.Error(
					w,
					http.StatusInternalServerError,
					"Failed to delete user",
				)
				return
			}

			_ = response.JSON(
				w,
				http.StatusOK,
				map[string]string{
					"message": "User deleted successfully",
				},
			)

		default:
			response.Error(
				w,
				http.StatusMethodNotAllowed,
				"Method not allowed",
			)
		}
	}
}
