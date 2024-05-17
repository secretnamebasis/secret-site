package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/api"
	"github.com/secretnamebasis/secret-site/app/middleware"
	"github.com/secretnamebasis/secret-site/app/views"
)

type resource struct {
	group fiber.Router
	name  string
	getAll,
	getByID,
	create,
	update,
	delete func(*fiber.Ctx) error
}

// Draw defines all the routes for the application
func Draw(app *fiber.App) {
	// Initialize middleware
	mw := middleware.New()
	app.Use(mw.LogRequests())

	// Define views routes
	defineViewsRoutes(app, mw)

	// Define API routes
	defineAPIRoutes(app, mw)
}

// defineViewsRoutes defines routes for views
func defineViewsRoutes(app *fiber.App, mw *middleware.Middleware) {

	// Create a route group for views
	viewsGroup := app.Group("/").Use(
		mw.HelmetMiddleware(),
		mw.RateLimiter(),
	)

	// Serve static files from the "assets" directory
	viewsGroup.Static("/", "./app/assets")
	viewsGroup.Static("/items", "./app/assets")

	// Define view routes using a map
	viewRoutes := map[string]func(*fiber.Ctx) error{
		"/":              views.Home,
		"/about":         views.About,
		"/items":         views.Items,
		"/items/new":     views.NewItem,
		"/items/:scid":   views.Item,
		"/images/:scid":  views.Images,
		"/files/:scid":   views.Files,
		"/users/":        views.Users,
		"/users/new":     views.NewUser,
		"/users/:wallet": views.User,
	}

	// Register view routes
	for path, handler := range viewRoutes {
		viewsGroup.Get(path, handler)
	}

	// Actions
	viewRoutes = map[string]func(*fiber.Ctx) error{
		"/users/submit": views.SubmitUser,
		"/items/submit": views.SubmitItem,
	}

	// Register post routes
	for path, handler := range viewRoutes {
		viewsGroup.Post(path, handler)
	}

}

// DefineAPIRoutes defines routes for APIs
func defineAPIRoutes(app *fiber.App, mw *middleware.Middleware) {
	var r = resource{}

	// Create a route group for API endpoints
	apiGroup := app.Group("/api")

	// Apply middleware for API endpoints
	apiGroup.Use(
		mw.HelmetMiddleware(),
		// would be nice to turn this on
		mw.RateLimiter(),
	)

	apiGroup.Get("/ping", api.Ping)

	// here there be monsters
	roles := []string{
		"user",
		// "admin", // we need to start thinking about user roles
	}

	apiGroup.Use(
		mw.AuthRequired(
			roles[0],
		),
	)

	r = resource{
		group:   apiGroup,
		name:    "items",
		getAll:  api.AllItems,
		getByID: api.ItemByID,
		create:  api.CreateItemOrder,
		update:  api.UpdateItem,
		delete:  api.DeleteItem,
	}

	// Define API routes for items
	defineItemRoutes(r)

	r = resource{
		group:   apiGroup,
		name:    "users",
		getAll:  api.AllUsers,
		getByID: api.UserByWallet,
		create:  api.CreateUserOrder,
		update:  api.UpdateUser,
		delete:  api.DeleteUser,
	}

	// Define API routes for users
	defineUserRoutes(r)
}

// Define resource routes for CRUD operations
func defineUserRoutes(r resource) {
	resource := r.group.Group("/" + r.name)
	resource.Get("/", r.getAll)
	resource.Post("/", r.create)
	resource.Get("/:wallet", r.getByID)
	resource.Put("/:wallet", r.update)
	resource.Delete("/:wallet", r.delete)
}

// Define resource routes for CRUD operations
func defineItemRoutes(r resource) {
	resource := r.group.Group("/" + r.name)
	resource.Get("/", r.getAll)
	resource.Post("/", r.create)
	resource.Get("/:scid", r.getByID)
	resource.Put("/:scid", r.update)
	resource.Delete("/:scid", r.delete)
}
