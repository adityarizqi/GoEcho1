package handler

import (
	"GoEcho1/db"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type DebugHandler struct {
	db *gorm.DB
}

func NewDebugHandler(gormDB *gorm.DB) *DebugHandler {
	return &DebugHandler{db: gormDB}
}

func (h *DebugHandler) ReseedDatabase(c echo.Context) error {
	err := db.ResetAndSeed(h.db)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to reseed database: "+err.Error())
	}
	return c.String(http.StatusOK, "Database has been successfully reset and re-seeded!")
}
