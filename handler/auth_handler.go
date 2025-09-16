package handler

import (
	"GoEcho1/service"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) ShowLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func (h *AuthHandler) Login(c echo.Context) error {
	nim := c.FormValue("nim")
	password := c.FormValue("password")

	token, user, err := h.authService.Login(nim, password)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/login?error=true")
	}

	// Set cookie
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)

	if user.Role == "dosen" {
		return c.Redirect(http.StatusSeeOther, "/dosen/dashboard")
	}
	return c.Redirect(http.StatusSeeOther, "/mahasiswa/dashboard")
}

func (h *AuthHandler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/login")
}
