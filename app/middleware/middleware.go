// middleware/middleware.go

package middleware

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/secretnamebasis/secret-site/app/api"
	"github.com/secretnamebasis/secret-site/app/database"
)

// Middleware provides a collection of middleware handlers
type Middleware struct{}

// New creates a new instance of Middleware
func New() *Middleware {
	return &Middleware{}
}

// instead of toggling these on and off, let's set up a "log-level"
func (m *Middleware) LogRequests() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Log request details
		log.Printf("Request: %s %s", c.Method(), c.OriginalURL())

		// Log request headers
		log.Println("Request Headers:")
		c.Request().Header.VisitAll(func(key, value []byte) {
			log.Printf("%s: %s", key, value)
		})

		// Log request body if present
		// this adds overhead to the processing of the server by 2x
		if len(c.Request().Body()) > 0 {
			log.Println("Request Body: " + string(c.Request().Body()))
		}

		// Proceed to next middleware or route handler
		if err := c.Next(); err != nil {
			return err
		}

		// Log response details
		// adds little overhead if any.
		log.Printf("Response: %d", c.Response().StatusCode())

		// // // Log response headers
		// // adds little overhead, but more noise
		log.Println("Response Headers:")
		c.Response().Header.VisitAll(func(key, value []byte) {
			log.Printf("%s: %s", key, value)
		})

		// // Log response body if present
		// this add trmendous insight, but causes the server to work 4x
		if len(c.Response().Body()) > 0 {
			log.Printf("Response Body: %s\n", string(c.Response().Body()))
		}

		return nil
	}
}

// AuthRequired middleware authenticates incoming requests and checks for required roles
func (m *Middleware) AuthRequired(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header from the request
		username, _, err := getCredentials(c)
		if err != nil {
			// Error getting credentials
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		_, err = database.GetUserByUsername(username)
		if err != nil {
			// Error extracting credentials, return unauthorized
			return api.ErrorResponse(c, fiber.StatusInternalServerError, "Unauthorized: we don't know you")
		}

		// if user == nil {
		// 	// User not found in the database, return unauthorized
		// 	return api.ErrorResponse(c, fiber.StatusInternalServerError, "Unauthorized: you don't exist")
		// }

		// Check if the user has any of the required roles
		// if !hasRole(user.Role, roles) {
		// 	return api.ErrorResponse(c, fiber.StatusInternalServerError, "Forbidden")
		// }

		// Proceed to the next middleware or route handler
		return c.Next()
	}
}
func getCredentials(c *fiber.Ctx) (username, password string, err error) {
	// Get the Authorization header from the request
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// No Authorization header found
		return "", "", errors.New("no Authorization header found")
	}

	// Extract the username and password from the Authorization header
	decodedCredentials, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
	if err != nil {
		// Error decoding credentials
		return "", "", fmt.Errorf("error decoding credentials: %v", err)
	}

	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) != 2 {
		// Invalid credentials format
		return "", "", errors.New("invalid credentials format")
	}

	return credentials[0], credentials[1], nil
}

// hasRole checks if the user has any of the required roles
func hasRole(userRoles []string, requiredRoles []string) bool {
	for _, role := range requiredRoles {
		for _, userRole := range userRoles {
			if role == userRole {
				return true
			}
		}
	}
	return false
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
