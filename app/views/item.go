package views

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/deroproject/derohe/rpc"
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

// ItemData defines the data structure for the item detail template
type ItemData struct {
	Title       string
	Address     string
	Item        models.Item
	SC_Data     rpc.GetSC_Result
	ImageUrl    string
	Image       string
	Description string
}

// Item renders the item detail page
func Item(c *fiber.Ctx) error {
	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Get the item ID from the request parameters
	scid := c.Params("scid")

	scData, err := dero.GetSCID(
		config.NodeEndpoint,
		scid,
	)
	if err != nil {
		return c.Status(
			fiber.StatusNotFound,
		).JSON(
			fiber.Map{
				"message": err.Error(),
				"status":  "error",
			},
		)
	}

	// Truncate key "C" value to display first 16 bytes followed by an ellipsis and then the next 16 bytes
	cValue := scData.VariableStringKeys["C"].(string)
	if len(cValue) > 32 {
		truncatedValue := cValue[:16] + "..." + cValue[len(cValue)-16:]
		scData.VariableStringKeys["C"] = truncatedValue
	}

	// Decode hex values in VariableStringKeys for all keys except "c"
	for k, v := range scData.VariableStringKeys {
		if k != "C" { // Exclude key "c"
			hexValue, ok := v.(string)
			if ok && isHex(hexValue) {
				decodedValue := decodeHex(hexValue)
				// Convert certain keys to string after decoding
				if k == "artificerAddr" || k == "creatorAddr" || k == "owner" {
					result, err := rpc.NewAddressFromCompressedKeys([]byte(decodedValue))
					if err != nil {
						return err
					}
					// Handle non-printable characters or characters that can't be directly represented as strings
					scData.VariableStringKeys[k] = result.String()
				} else {
					scData.VariableStringKeys[k] = decodedValue
				}
			}
		}
	}

	// Retrieve the item by ID
	item, err := controllers.GetItemBySCID(scid)
	if err != nil {
		return c.Status(
			fiber.StatusNotFound,
		).JSON(
			fiber.Map{
				"message": err.Error(),
				"status":  "error",
			},
		)
	}

	var itemData models.ItemData
	if err := json.Unmarshal(item.Data, &itemData); err != nil {
		return c.Status(
			fiber.StatusNotFound,
		).JSON(
			fiber.Map{
				"message": err.Error(),
				"status":  "error",
			},
		)
	}
	// Define data for rendering the template
	data := ItemData{
		Title:       config.Domain,
		Address:     addr.String(),
		Item:        item,
		SC_Data:     *scData,
		ImageUrl:    item.ImageURL,
		Image:       itemData.Image,
		Description: itemData.Description,
	}

	// Render the template using renderTemplate function
	if err := renderTemplate(c, "app/public/item.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

// isHex checks if a string is in hexadecimal format
func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// decodeHex decodes a hexadecimal string to its ASCII representation
func decodeHex(hexString string) string {
	var str strings.Builder
	for i := 0; i < len(hexString); i += 2 {
		b, _ := strconv.ParseUint(hexString[i:i+2], 16, 8)
		str.WriteByte(byte(b))
	}
	return str.String()
}
