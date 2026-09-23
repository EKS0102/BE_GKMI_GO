package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

			query := r.URL.Query()

			search := strings.TrimSpace(query.Get("search"))
			jenisKelamin := strings.TrimSpace(query.Get("jenis_kelamin"))
			statusJemaat := strings.TrimSpace(query.Get("status_jemaat"))
			statusDiakonia := strings.TrimSpace(query.Get("status_diakonia"))
			kelompokIbadah := strings.TrimSpace(query.Get("kelompok_ibadah"))

			sortBy := strings.TrimSpace(query.Get("sort_by"))
			sortOrder := strings.TrimSpace(query.Get("sort_order"))

			if sortBy == "" {
				sortBy = "id"
			}

			if sortOrder == "" {
				sortOrder = "asc"
			}

			pageParam := strings.TrimSpace(query.Get("page"))
			limitParam := strings.TrimSpace(query.Get("limit"))

			hasFilter := search != "" ||
				jenisKelamin != "" ||
				statusJemaat != "" ||
				statusDiakonia != "" ||
				kelompokIbadah != "" ||
				query.Get("sort_by") != "" ||
				query.Get("sort_order") != ""

			hasPagination := pageParam != "" || limitParam != ""

			if !hasFilter && !hasPagination {

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

			if pageParam != "" {
				value, err := strconv.Atoi(pageParam)
				if err != nil {
					response.Error(
						w,
						http.StatusBadRequest,
						"Invalid page",
					)
					return
				}

				page = value
			}

			if limitParam != "" {
				value, err := strconv.Atoi(limitParam)
				if err != nil {
					response.Error(
						w,
						http.StatusBadRequest,
						"Invalid limit",
					)
					return
				}

				limit = value
			}

			if !hasFilter {
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
				return
			}

			result, err := jemaatService.GetFilteredPaginated(
				r.Context(),
				search,
				jenisKelamin,
				statusJemaat,
				statusDiakonia,
				kelompokIbadah,
				sortBy,
				sortOrder,
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

			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				response.Error(
					w,
					http.StatusBadRequest,
					"Invalid JSON body",
				)
				return
			}

			if message := request.Validate(); message != "" {
				response.Error(
					w,
					http.StatusBadRequest,
					message,
				)
				return
			}

			jemaat, err := jemaatService.Create(
				r.Context(),
				request,
			)

			if err != nil {
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

			idString := strings.TrimPrefix(
				r.URL.Path,
				"/api/jemaat/",
			)

			id, err := strconv.Atoi(idString)
			if err != nil {
				response.Error(
					w,
					http.StatusBadRequest,
					"Invalid jemaat ID",
				)
				return
			}

			var request model.UpdateJemaatRequest

			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				response.Error(
					w,
					http.StatusBadRequest,
					"Invalid JSON body",
				)
				return
			}

			if message := request.Validate(); message != "" {
				response.Error(
					w,
					http.StatusBadRequest,
					message,
				)
				return
			}

			jemaat, err := jemaatService.Update(
				r.Context(),
				id,
				request,
			)

			if err != nil {
				if strings.Contains(
					strings.ToLower(err.Error()),
					"not found",
				) {
					response.Error(
						w,
						http.StatusNotFound,
						"Jemaat not found",
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

			idString := strings.TrimPrefix(
				r.URL.Path,
				"/api/jemaat/",
			)

			id, err := strconv.Atoi(idString)
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
				if strings.Contains(
					strings.ToLower(err.Error()),
					"not found",
				) {
					response.Error(
						w,
						http.StatusNotFound,
						"Jemaat not found",
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
