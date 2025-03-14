package util

import (
	"net/http"

	response "github.com/federicodosantos/image-smith/pkg/response"
	"github.com/jmoiron/sqlx"
)

func HealthCheck(router *http.ServeMux, db *sqlx.DB) {
	type HealthStatus struct {
		Status   string `json:"status"`
		Database string `json:"database"`
	}

	router.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		status := HealthStatus{
			Status:   "healthy",
			Database: "healthy",
		}

		if err := db.Ping(); err != nil {
			status.Status = "unhealthy"
			status.Database = "unhealthy"
		}

		httpStatus := http.StatusOK
		if status.Status != "healthy" {
			httpStatus = http.StatusServiceUnavailable
		}

		response.SuccessResponse(w, httpStatus, "health check", status)

	})
}
