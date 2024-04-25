package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/secretnamebasis/secret-site/app/api"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAPIInfo(t *testing.T) {
	t.Run("Given APIInfo handler", func(t *testing.T) {
		app := fiber.New()
		app.Get("/ping", api.Ping)

		t.Run("When requesting APIInfo", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			resp, err := app.Test(req)

			t.Run("Then return APIInfo successfully", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				assert.NoError(t, err)

				assert.Equal(t, "pong", responseBody["data"])
				assert.Equal(t, "success", responseBody["status"])
			})
		})
	})
}
