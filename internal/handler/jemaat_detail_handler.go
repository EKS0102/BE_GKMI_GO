package handler

import (
	"net/http"
	"strconv"
	"strings"

	"BE_GKMI_NTC_GO/internal/response"
)

func GetJemaatByID(jemaatService JemaatServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := strings.TrimPrefix(r.URL.Path, "/api/jemaat/")
		id, err := strconv.Atoi(idString)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid jemaat ID")
			return
		}

		jemaat, err := jemaatService.GetByID(r.Context(), id)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "not found") {
				response.Error(w, http.StatusNotFound, "Jemaat not found")
				return
			}
			response.Error(w, http.StatusInternalServerError, "Failed to get jemaat")
			return
		}

		response.JSON(w, http.StatusOK, jemaat)
	}
}
