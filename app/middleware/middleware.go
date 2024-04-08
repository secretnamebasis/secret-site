// middleware/middleware.go

package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/secretnamebasis/secret-site/app/controllers"
)

// Middleware provides a collection of middleware handlers
type Middleware struct{}

// New creates a new instance of Middleware
func New() *Middleware {
	return &Middleware{}
}

func (m *Middleware) LogRequests() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log request details
		log.Printf("Request: %s %s", c.Method(), c.OriginalURL())

		// Log request headers
		log.Println("Request Headers:")
		c.Request().Header.VisitAll(func(key, value []byte) {
			log.Printf("%s: %s", key, value)
		})

		// // Log request body if present
		if len(c.Request().Body()) > 0 {
			log.Println("Request Body:")
			log.Println(string(c.Request().Body()))
		}

		// Proceed to next middleware or route handler
		if err := c.Next(); err != nil {
			return err
		}

		// Log response details
		log.Printf("Response: %d", c.Response().StatusCode())

		// // Log response headers
		log.Println("Response Headers:")
		c.Response().Header.VisitAll(func(key, value []byte) {
			log.Printf("%s: %s", key, value)
		})

		// Log response body if present
		if len(c.Response().Body()) > 0 {
			log.Printf("Response Body: %s\n", string(c.Response().Body()))
		}

		return nil
	}
}

// AuthRequired middleware authenticates incoming requests
func (m *Middleware) AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := controllers.AllUsers(c)
		if err != nil {
			log.Println("Error fetching users:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		userMap := make(map[string]string)
		for _, user := range users {
			userMap[user.User] = user.Password
		}

		cfg := basicauth.Config{
			Users: userMap,
		}

		authMiddleware := basicauth.New(cfg)
		if authMiddleware == nil {
			log.Println("Error creating basic auth middleware")
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		return authMiddleware(c)
	}
}

// RateLimiter middleware limits the rate of incoming requests
func (m *Middleware) RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,             // Maximum number of requests allowed in Expiration duration
		Expiration: 1 * time.Minute, // Time duration for which requests are tracked
		KeyGenerator: func(c *fiber.Ctx) string { // Generate a key for identifying requests
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error { // Handler when limit is reached
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests")
		},
	})
}

// HelmetMiddleware returns a middleware handler for Helmet security
func (m *Middleware) HelmetMiddleware() fiber.Handler {
	return helmet.New()
}
