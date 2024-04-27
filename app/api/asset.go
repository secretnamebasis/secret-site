package api

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/models"
)

func CreateAssetOrder(c *fiber.Ctx) error {
	var order models.JSON_Asset_Order
	form, _ := c.MultipartForm()
	// if err != nil {
	// 	return err
	// }
	if form != nil {
		if err := processAssetOrderForm(
			form,
			&order,
		); err != nil {
			return ErrorResponse(
				c,
				fiber.StatusBadRequest,
				err.Error(),
			)
		}
	} else {
		// Parse request body into new item
		if err := c.BodyParser(&order); err != nil {
			return ErrorResponse(
				c,
				fiber.StatusBadRequest,
				err.Error(),
			)
		}

	}
	// currently the order doesn't need to have user pass
	// we will be able to
	// if err := processAssetOrderCredentials(
	// 	c,
	// 	&order,
	// ); err != nil {
	// 	return ErrorResponse(
	// 		c,
	// 		fiber.StatusInternalServerError,
	// 		err.Error(),
	// 	)
	// }
	// we need to make a controller
	checkout, err := controllers.CreateAssetCheckout(&order)
	if err != nil {
		return ErrorResponse(
			c,
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}

	// Return success response
	return SuccessResponse(
		c,
		"address expires at: "+checkout.Expiration.String(),
		checkout.Address,
	)
}

// for now, I don't know how useful it would be to process creds unless there
// was some kind of subscription basis or if they were a secret member
// the point is, we have this off at the moment because anyone should be able
// to claim an asset by paying for it. Now we need to extract
// dero addresses and process them async
// func processAssetOrderCredentials(c *fiber.Ctx, order *models.JSON_Asset_Order) error {
// 	if err := controllers.ValidateWalletAddress(order.Wallet); err != nil {
// 		return ErrorResponse(
// 			c,
// 			fiber.StatusInternalServerError,
// 			err.Error(),
// 		)
// 	}

// 	name, pass, err := getCredentials(c)
// 	if err != nil {
// 		return ErrorResponse(
// 			c,
// 			fiber.StatusInternalServerError,
// 			err.Error(),
// 		)
// 	}

// 	if name == "" || pass == "" {

// 		return ErrorResponse(
// 			c,
// 			fiber.StatusInternalServerError,
// 			"name or pass is empty",
// 		)
// 	}

//		return nil
//	}
func processAssetOrderForm(form *multipart.Form, order *models.JSON_Asset_Order) error {
	order.Name = form.Value["name"][0]
	order.Description = form.Value["description"][0]
	order.Type = form.Value["type"][0]
	order.Wallet = form.Value["wallet"][0]
	order.Collection = form.Value["collection"][0]
	return nil
}
