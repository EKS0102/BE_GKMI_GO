package handler

import (
	"net/http"

	"BE_GKMI_NTC_GO/internal/model"
	"BE_GKMI_NTC_GO/internal/response"
)

// Health godoc
// @Summary Health check
// @Description Check whether the backend API is running
// @Tags Health
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Failure 405 {object} map[string]string
// @Router /api/health [get]
func Health(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		response.Error(
			w,
			http.StatusMethodNotAllowed,
			"Method not allowed",
		)
		return
	}

	data := model.HealthResponse{
		Status:      "ok",
		Application: "BE_GKMI_NTC_GO",
		Version:     "1.0.0",
	}

	_ = response.JSON(
		w,
		http.StatusOK,
		data,
	)
}
