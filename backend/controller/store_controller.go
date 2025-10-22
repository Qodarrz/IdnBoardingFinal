// controller/store.go
package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	dto "github.com/Qodarrz/fiber-app/dto"
	helpers "github.com/Qodarrz/fiber-app/helper"
	"github.com/Qodarrz/fiber-app/middleware"
	service "github.com/Qodarrz/fiber-app/service"
	"github.com/gofiber/fiber/v2"
)

type StoreController struct {
	storeService service.StoreServiceInterface
}

func InitStoreController(app *fiber.App, svc service.StoreServiceInterface, mw *middleware.Middlewares) {
	ctrl := &StoreController{storeService: svc}

	public := app.Group("/api/store")
	public.Get("/items", ctrl.GetAllStoreItems)
	public.Get("/items/:id", ctrl.GetStoreItemByID)

	private := app.Group("/api/store", mw.JWT)
	private.Post("/items", ctrl.CreateStoreItem)
	private.Put("/items/:id", ctrl.UpdateStoreItem)
	private.Delete("/items/:id", ctrl.DeleteStoreItem)
	private.Post("/orders/:item_id", ctrl.CreateOrder)
	private.Get("/orders", ctrl.GetUserOrders)
	private.Get("/orders/:id", ctrl.GetOrderByID)
	private.Post("/orders/:id/cancel", ctrl.CancelOrder)
	// private.Post("/orders/item/:item_id", ctrl.CreateOrderByItemID)
}


func (c *StoreController) GetAllStoreItems(ctx *fiber.Ctx) error {
	status := ctx.Query("status", "")

	items, err := c.storeService.GetAllStoreItems(ctx.Context(), status)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "store items retrieved successfully", items))
}

func (c *StoreController) GetStoreItemByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid item ID"))
	}

	item, err := c.storeService.GetStoreItemByID(ctx.Context(), id)
	if err != nil {
		if err.Error() == "store item not found" {
			return ctx.Status(http.StatusNotFound).JSON(helpers.BasicResponse(false, err.Error()))
		}
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "store item retrieved successfully", item))
}

func (c *StoreController) CreateStoreItem(ctx *fiber.Ctx) error {
    req := new(dto.CreateStoreItemDTO)

    req.Name = ctx.FormValue("name")
    req.Description = ctx.FormValue("description")
    req.PricePoints, _ = strconv.Atoi(ctx.FormValue("price_points"))
    req.Stock, _ = strconv.Atoi(ctx.FormValue("stock"))

    fileHeader, err := ctx.FormFile("image")
    if err == nil && fileHeader != nil {
        tempPath := filepath.Join(os.TempDir(), fileHeader.Filename)
        if err := ctx.SaveFile(fileHeader, tempPath); err != nil {
            return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": err.Error(),
            })
        }
        url, err := helpers.UploadFile(tempPath)
        if err != nil {
            return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": err.Error(),
            })
        }
        req.ImageURL = url
    }

    // Validasi manual
    if req.Name == "" || req.PricePoints <= 0 {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "name wajib diisi dan price_points harus > 0",
        })
    }

    item, err := c.storeService.CreateStoreItem(ctx.Context(), req)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": err.Error(),
        })
    }

    return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "Store item created successfully",
        "data":    item,
    })
}


func (c *StoreController) UpdateStoreItem(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid item ID"))
	}

	req := new(dto.UpdateStoreItemDTO)
	if err := helpers.BindAndValidate(ctx, req); err != nil {
		if vErr, ok := err.(*helpers.ValidationError); ok {
			return ctx.Status(http.StatusBadRequest).JSON(helpers.ErrorResponseRequest(false, vErr.Message, vErr.Errors))
		}
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	item, err := c.storeService.UpdateStoreItem(ctx.Context(), id, req)
	if err != nil {
		if err.Error() == "store item not found" {
			return ctx.Status(http.StatusNotFound).JSON(helpers.BasicResponse(false, err.Error()))
		}
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "store item updated successfully", item))
}

func (c *StoreController) DeleteStoreItem(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid item ID"))
	}

	err = c.storeService.DeleteStoreItem(ctx.Context(), id)
	if err != nil {
		if err.Error() == "store item not found" {
			return ctx.Status(http.StatusNotFound).JSON(helpers.BasicResponse(false, err.Error()))
		}
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.BasicResponse(true, "store item deleted successfully"))
}

func (c *StoreController) CreateOrder(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(
			helpers.BasicResponse(false, "invalid token"),
		)
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			helpers.BasicResponse(false, "invalid user ID"),
		)
	}

	itemID, err := strconv.ParseInt(ctx.Params("item_id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			helpers.BasicResponse(false, "invalid item_id"),
		)
	}

	qty, err := strconv.Atoi(ctx.Query("qty", "1"))
	if err != nil || qty <= 0 {
		qty = 1
	}

	// ✅ gunakan CreateOrderByItemID
	orderResponse, err := c.storeService.CreateOrderByItemID(ctx.Context(), userID, itemID, qty)
	if err != nil {
		switch err.Error() {
		case "insufficient points", "item not found", "insufficient stock", "item is not available":
			return ctx.Status(http.StatusBadRequest).JSON(
				helpers.BasicResponse(false, err.Error()),
			)
		default:
			return ctx.Status(http.StatusInternalServerError).JSON(
				helpers.BasicResponse(false, err.Error()),
			)
		}
	}

	return ctx.Status(http.StatusCreated).JSON(
		helpers.SuccessResponseWithData(true, "order created successfully", orderResponse),
	)
}



func (c *StoreController) GetOrderByID(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "invalid token"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid user ID"))
	}

	orderID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid order ID"))
	}

	order, err := c.storeService.GetOrderByID(ctx.Context(), orderID)
	if err != nil {
		if err.Error() == "order not found" {
			return ctx.Status(http.StatusNotFound).JSON(helpers.BasicResponse(false, err.Error()))
		}
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	// Check if the order belongs to the authenticated user
	if order.UserID != userID {
		return ctx.Status(http.StatusForbidden).JSON(helpers.BasicResponse(false, "access denied"))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "order retrieved successfully", order))
}

func (c *StoreController) GetUserOrders(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "invalid token"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid user ID"))
	}

	orders, err := c.storeService.GetUserOrders(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "user orders retrieved successfully", orders))
}

func (c *StoreController) CancelOrder(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "invalid token"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid user ID"))
	}

	orderID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid order ID"))
	}

	err = c.storeService.CancelOrder(ctx.Context(), userID, orderID)
	if err != nil {
		switch err.Error() {
		case "order not found", "unauthorized to cancel this order", "only pending orders can be cancelled":
			return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
		default:
			return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
		}
	}

	return ctx.Status(http.StatusOK).JSON(helpers.BasicResponse(true, "order cancelled successfully"))
}