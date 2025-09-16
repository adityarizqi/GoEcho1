package main

import (
	"GoEcho1/db"
	"GoEcho1/handler"
	"GoEcho1/middleware"
	"GoEcho1/repository"
	"GoEcho1/service"
	"html/template"
	"io"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	// Initialize Database
	gormDB := db.InitDB()
	db.Seed(gormDB)

	// Dependency Injection
	userRepo := repository.NewUserRepository(gormDB)
	courseRepo := repository.NewCourseRepository(gormDB)
	krsRepo := repository.NewKRSRepository(gormDB)

	authService := service.NewAuthService(userRepo)
	krsService := service.NewKRSService(krsRepo, courseRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService)
	mhsHandler := handler.NewMahasiswaHandler(krsService)
	dosenHandler := handler.NewDosenHandler(krsService)
	debugHandler := handler.NewDebugHandler(gormDB)

	// Echo instance
	e := echo.New()

	// Setup template renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// --- Routes ---
	// Public routes
	e.GET("/login", authHandler.ShowLoginPage)
	e.POST("/login", authHandler.Login)
	e.GET("/logout", authHandler.Logout)
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(302, "/login")
	})

	// Mahasiswa routes
	mhsGroup := e.Group("/mahasiswa")
	mhsGroup.Use(middleware.RoleMiddleware("mahasiswa"))
	mhsGroup.GET("/dashboard", mhsHandler.ShowDashboard)
	mhsGroup.POST("/enroll", mhsHandler.EnrollCourse)

	// Dosen routes
	dosenGroup := e.Group("/dosen")
	dosenGroup.Use(middleware.RoleMiddleware("dosen"))
	dosenGroup.GET("/dashboard", dosenHandler.ShowDashboard)

	// Debug routes
	debugGroup := e.Group("/debug")
	debugGroup.POST("/reseed", debugHandler.ReseedDatabase)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
