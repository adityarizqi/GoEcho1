package handler

import (
	"GoEcho1/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DosenHandler struct {
	krsService service.KRSService
}

func NewDosenHandler(krsService service.KRSService) *DosenHandler {
	return &DosenHandler{krsService}
}

func (h *DosenHandler) ShowDashboard(c echo.Context) error {
	data, err := h.krsService.GetDashboardDataForDosen()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Gagal memuat data")
	}

	return c.Render(http.StatusOK, "dosen_dashboard.html", data)
}
