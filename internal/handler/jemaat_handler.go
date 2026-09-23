package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/repository"
	"BE_GKMI_NTC_GO/internal/response"
	"BE_GKMI_NTC_GO/internal/service"
)

func Jemaat(
	jemaatRepository *repository.JemaatRepository,
	jemaatService *service.JemaatService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			jemaatList, err := jemaatService.GetAll(r.Context())
			if err != nil {
				response.Error(w, http.StatusInternalServerError, "Failed to get jemaat")
				return
			}

			_ = response.JSON(w, http.StatusOK, jemaatList)

		case http.MethodPost:
			var request model.CreateJemaatRequest

			err := json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				response.Error(
					w,
					http.StatusBadRequest,
					"Invalid request body",
				)
				return
			}

			jemaat, err := jemaatService.Create(
				r.Context(),
				request,
			)
			if err != nil {
				if validationErr, ok := err.(*service.ValidationError); ok {
					response.Error(
						w,
						http.StatusBadRequest,
						validationErr.Message,
					)
					return
				}

				response.Error(
					w,
					http.StatusInternalServerError,
					"Failed to create jemaat",
				)
				return
			}

			_ = response.JSON(
				w,
				http.StatusCreated,
				jemaat,
			)

		case http.MethodPut:
			idStr := r.URL.Path[len("/api/jemaat/"):]

			id, err := strconv.Atoi(idStr)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Invalid jemaat ID")
				return
			}

			var request model.UpdateJemaatRequest

			err = json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Invalid request body")
				return
			}

			jemaat, err := jemaatService.Update(
				r.Context(),
				id,
				request,
			)
			if err != nil {
				if validationErr, ok := err.(*service.ValidationError); ok {
					response.Error(
						w,
						http.StatusBadRequest,
						validationErr.Message,
					)
					return
				}

				response.Error(
					w,
					http.StatusInternalServerError,
					"Failed to update jemaat",
				)
				return
			}

			_ = response.JSON(
				w,
				http.StatusOK,
				jemaat,
			)

		case http.MethodDelete:
			idStr := r.URL.Path[len("/api/jemaat/"):]

			id, err := strconv.Atoi(idStr)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "Invalid jemaat ID")
				return
			}

			err = jemaatService.Delete(r.Context(), id)
			if err != nil {
				response.Error(
					w,
					http.StatusInternalServerError,
					"Failed to delete jemaat",
				)
				return
			}

			_ = response.JSON(
				w,
				http.StatusOK,
				map[string]string{
					"message": "Jemaat deleted successfully",
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
