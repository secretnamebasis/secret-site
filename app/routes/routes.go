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

	// Define view routes
	viewRoutes := []struct {
		Path   string
		Handle func(*fiber.Ctx) error
	}{
		{
			Path:   "/",
			Handle: views.Home,
		},
		{
			Path:   "/about",
			Handle: views.About,
		},
		{
			Path:   "/items",
			Handle: views.Items,
		},
		{
			Path:   "/items/new",
			Handle: views.NewItem,
		},
		{
			Path:   "/items/:scid",
			Handle: views.Item,
		},
		{
			Path:   "/images/:scid",
			Handle: views.Images,
		},
		{
			Path:   "/files/:scid",
			Handle: views.Files,
		},
		{
			Path:   "/users/",
			Handle: views.Users,
		},
		{
			Path:   "/users/new",
			Handle: views.NewUser,
		},
		{
			Path:   "/users/:wallet",
			Handle: views.User,
		},
	}

	// Register view routes
	for _, route := range viewRoutes {
		viewsGroup.Get(
			route.Path,
			route.Handle,
		)
	}
	// Actions
	viewsGroup.Post("/users/submit", views.SubmitUser)
	viewsGroup.Post("/items/submit", views.SubmitItem)

}

// DefineAPIRoutes defines routes for APIs
func defineAPIRoutes(app *fiber.App, mw *middleware.Middleware) {

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
	roles := []string{"user"}
	apiGroup.Use(mw.AuthRequired(roles[0]))

	// Define API routes for items
	defineResourceRoutes(
		apiGroup,
		"items",
		api.AllItems,
		api.ItemByID,
		api.CreateItemOrder,
		api.UpdateItem,
		api.DeleteItem,
	)

	// Define API routes for users
	defineResourceRoutes(
		apiGroup,
		"users",
		api.AllUsers,
		api.UserByID,
		api.CreateUserOrder,
		api.UpdateUser,
		api.DeleteUser,
	)
}

// Define resource routes for CRUD operations
func defineResourceRoutes(
	group fiber.Router,
	resourceName string,
	getAll,
	getByID,
	create,
	update,
	delete func(*fiber.Ctx) error,
) {
	resource := group.Group("/" + resourceName)
	resource.Get("/", getAll)
	resource.Post("/", create)
	resource.Get("/:id", getByID)
	resource.Put("/:id", update)
	resource.Delete("/:id", delete)
}
