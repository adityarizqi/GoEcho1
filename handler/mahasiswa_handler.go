package handler

import (
	"GoEcho1/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MahasiswaHandler struct {
	krsService service.KRSService
}

func NewMahasiswaHandler(krsService service.KRSService) *MahasiswaHandler {
	return &MahasiswaHandler{krsService}
}

func (h *MahasiswaHandler) ShowDashboard(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	data, err := h.krsService.GetDashboardDataForStudent(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Gagal memuat data")
	}
	return c.Render(http.StatusOK, "mahasiswa_dashboard.html", data)
}

func (h *MahasiswaHandler) EnrollCourse(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	courseIDStr := c.FormValue("course_id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Course ID tidak valid")
	}

	err = h.krsService.EnrollCourse(userID, uint(courseID))
	if err != nil {
		// Di dunia nyata, Anda akan menampilkan error ini dengan lebih baik
		return c.String(http.StatusConflict, err.Error())
	}

	return c.Redirect(http.StatusSeeOther, "/mahasiswa/dashboard")
}
