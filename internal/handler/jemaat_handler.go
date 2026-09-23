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
			pageStr := r.URL.Query().Get("page")
			limitStr := r.URL.Query().Get("limit")

			if pageStr == "" && limitStr == "" {
				jemaatList, err := jemaatService.GetAll(r.Context())
				if err != nil {
					response.Error(
						w,
						http.StatusInternalServerError,
						"Failed to get jemaat",
					)
					return
				}

				_ = response.JSON(
					w,
					http.StatusOK,
					jemaatList,
				)
				return
			}

			page := 1
			limit := 10

			var err error

			if pageStr != "" {
				page, err = strconv.Atoi(pageStr)
				if err != nil {
					response.Error(
						w,
						http.StatusBadRequest,
						"Invalid page",
					)
					return
				}
			}

			if limitStr != "" {
				limit, err = strconv.Atoi(limitStr)
				if err != nil {
					response.Error(
						w,
						http.StatusBadRequest,
						"Invalid limit",
					)
					return
				}
			}

			result, err := jemaatService.GetPaginated(
				r.Context(),
				page,
				limit,
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
					"Failed to get jemaat",
				)
				return
			}

			_ = response.JSON(
				w,
				http.StatusOK,
				result,
			)

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
				response.Error(
					w,
					http.StatusBadRequest,
					"Invalid jemaat ID",
				)
				return
			}

			var request model.UpdateJemaatRequest

			err = json.NewDecoder(r.Body).Decode(&request)
			if err != nil {
				response.Error(
					w,
					http.StatusBadRequest,
					"Invalid request body",
				)
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

				if notFoundErr, ok := err.(*service.NotFoundError); ok {
					response.Error(
						w,
						http.StatusNotFound,
						notFoundErr.Message,
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
				response.Error(
					w,
					http.StatusBadRequest,
					"Invalid jemaat ID",
				)
				return
			}

			err = jemaatService.Delete(
				r.Context(),
				id,
			)
			if err != nil {
				if notFoundErr, ok := err.(*service.NotFoundError); ok {
					response.Error(
						w,
						http.StatusNotFound,
						notFoundErr.Message,
					)
					return
				}

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
