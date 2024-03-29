package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/api"
	"github.com/secretnamebasis/secret-site/app/middleware"
	"github.com/secretnamebasis/secret-site/app/views"
)

// Draw defines all the routes for the application
func Draw(app *fiber.App) {
	// Initialize middleware
	mw := middleware.New()

	// Define views routes
	defineViewsRoutes(app, mw)

	// Define API routes
	defineAPIRoutes(app, mw)
}

// defineViewsRoutes defines routes for views
func defineViewsRoutes(app *fiber.App, mw *middleware.Middleware) {
	// Create a route group for views
	viewsGroup := app.Group("/")
	viewsGroup.Use(mw.HelmetMiddleware(), mw.LogRequests(), mw.RateLimiter())

	// Define view routes
	viewsGroup.Get("/", views.Home)
	viewsGroup.Get("/items/:id", views.Item)
	viewsGroup.Get("/api/ping", api.Ping)
}

// defineAPIRoutes defines routes for APIs
func defineAPIRoutes(app *fiber.App, mw *middleware.Middleware) {
	// Create a route group for API endpoints
	apiGroup := app.Group("/api")

	// Apply middleware for API endpoints
	apiGroup.Use(mw.HelmetMiddleware(), mw.AuthRequired(), mw.LogRequests(), mw.RateLimiter())

	// Define API routes for items
	items := apiGroup.Group("/items")
	items.Get("/", api.AllItems)
	items.Get("/:id", api.ItemByID)
	items.Post("/", api.CreateItem)
	items.Put("/:id", api.UpdateItem)
	items.Delete("/:id", api.DeleteItem)

	// Define API routes for users
	users := apiGroup.Group("/users")
	users.Get("/", api.AllUsers)
	users.Get("/:id", api.UserByID)
	users.Post("/", api.CreateUser)
	users.Put("/:id", api.UpdateUser)
	users.Delete("/:id", api.DeleteUser)
}
