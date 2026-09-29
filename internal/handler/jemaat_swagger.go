package handler

// JemaatListSwaggerDoc documents GET /api/jemaat.
// @Summary List jemaat
// @Description Get jemaat list with optional search, filters, sorting, and pagination.
// @Tags Jemaat
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search by nama panggilan or nama lengkap"
// @Param jenis_kelamin query string false "Jenis kelamin" Enums(Laki-Laki,Perempuan)
// @Param status_jemaat query string false "Status jemaat" Enums(Jemaat,Simpatisan,Tamu)
// @Param status_diakonia query string false "Status diakonia" Enums(Ya,Tidak)
// @Param kelompok_ibadah query string false "Kelompok ibadah" Enums(Sekolah Minggu,Youth,Kompak,Koper,Lansia)
// @Param sort_by query string false "Sort field" Enums(id,nama_panggilan,nama_lengkap,tanggal_lahir,jenis_kelamin,status_jemaat,status_diakonia,kelompok_ibadah)
// @Param sort_order query string false "Sort order" Enums(asc,desc)
// @Param page query int false "Page number" minimum(1)
// @Param limit query int false "Items per page" minimum(1)
// @Success 200 {object} map[string]interface{} "Jemaat list or paginated response"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/jemaat [get]
func JemaatListSwaggerDoc() {}

// JemaatCreateSwaggerDoc documents POST /api/jemaat.
// @Summary Create jemaat
// @Tags Jemaat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateJemaatRequest true "Jemaat data"
// @Success 201 {object} model.Jemaat
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/jemaat [post]
func JemaatCreateSwaggerDoc() {}

// JemaatBulkSwaggerDoc documents POST /api/jemaat/bulk.
// @Summary Bulk create jemaat
// @Tags Jemaat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body []model.CreateJemaatRequest true "Jemaat data list"
// @Success 201 {array} model.Jemaat
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/jemaat/bulk [post]
func JemaatBulkSwaggerDoc() {}

// JemaatUpdateSwaggerDoc documents PUT /api/jemaat/{id}.
// @Summary Update jemaat
// @Tags Jemaat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Jemaat ID"
// @Param request body model.UpdateJemaatRequest true "Jemaat data"
// @Success 200 {object} model.Jemaat
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/jemaat/{id} [put]
func JemaatUpdateSwaggerDoc() {}

// JemaatDeleteSwaggerDoc documents DELETE /api/jemaat/{id}.
// @Summary Delete jemaat
// @Tags Jemaat
// @Produce json
// @Security BearerAuth
// @Param id path int true "Jemaat ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/jemaat/{id} [delete]
func JemaatDeleteSwaggerDoc() {}

// JemaatDetailSwaggerDoc documents GET /api/jemaat/{id}.
// @Summary Get jemaat by ID
// @Tags Jemaat
// @Produce json
// @Security BearerAuth
// @Param id path int true "Jemaat ID"
// @Success 200 {object} model.Jemaat
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/jemaat/{id} [get]
func JemaatDetailSwaggerDoc() {}
