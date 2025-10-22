// controller/carbon_controller.go
package controller

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	dto "github.com/Qodarrz/fiber-app/dto"
	helpers "github.com/Qodarrz/fiber-app/helper"
	"github.com/Qodarrz/fiber-app/middleware"
	service "github.com/Qodarrz/fiber-app/service"
	"github.com/gofiber/fiber/v2"
)

type CarbonController struct {
	carbonService service.CarbonServiceInterface
}

func InitCarbonController(app *fiber.App, svc service.CarbonServiceInterface, mw *middleware.Middlewares) {
	ctrl := &CarbonController{carbonService: svc}

	public := app.Group("/api/carbon", mw.JWT)
	
	public.Post("/vehicle", ctrl.CreateVehicle)
	public.Get("/vehicles", ctrl.ListUserVehicles)	
	public.Delete("/vehicle/:id", ctrl.DeleteVehicle)
	public.Patch("/vehicle/:id", ctrl.EditVehicle)
	public.Post("/vehicle-log", ctrl.AddVehicleLog)
	public.Get("/vehicle/:id/logs", ctrl.GetVehicleLogs)
	public.Get("/vehicle/logs", ctrl.GetAllVehicleLogs)
	public.Get("/vehicle/logs/:id", ctrl.GetVehicleLogByID)

	public.Post("/electronic", ctrl.CreateElectronic)
	public.Get("/electronics", ctrl.ListUserElectronics)
	public.Patch("/electronics/:id", ctrl.EditElectronic)
	public.Delete("/electronics/:id", ctrl.DeleteElectronic)
	public.Post("/electronics-log", ctrl.AddElectronicsLog)
	public.Get("/electronic/:id/logs", ctrl.GetElectronicsLogs)
	public.Get("/electronics/logs", ctrl.GetAllElectronicLogs)
	
}

func (c *CarbonController) CreateVehicle(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	dto := new(dto.CreateVehicleDTO)
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid body"))
	}

	vehicle, err := c.carbonService.CreateVehicle(ctx.Context(), userID, dto)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(helpers.SuccessResponseWithData(true, "vehicle created successfully", vehicle))
}

func (c *CarbonController) ListUserVehicles(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	vehicles, err := c.carbonService.ListUserVehicles(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "vehicles retrieved successfully", vehicles))
}

func (c *CarbonController) GetVehicleLogByID(ctx *fiber.Ctx) error {
    // Ambil userID dari claims
    claims := helpers.GetUserClaims(ctx)
    userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

    // Ambil log ID dari params
    logID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid log ID"))
    }
	
    // Panggil service untuk ambil log per logID + cek ownership userID
    log, err := c.carbonService.GetVehicleLogByID(ctx.Context(), userID, logID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return ctx.Status(fiber.StatusNotFound).JSON(helpers.BasicResponse(false, "log not found"))
        }
        return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
    }

    return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "vehicle log retrieved successfully", log))
}



func (c *CarbonController) AddVehicleLog(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	dto := new(dto.AddVehicleLogDTO)
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid body"))
	}

	if err := c.carbonService.AddVehicleLog(ctx.Context(), userID, dto); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.BasicResponse(true, "vehicle log berhasil ditambahkan"))
}

func (c *CarbonController) GetVehicleLogs(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	vehicleID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid vehicle ID"))
	}

	logs, err := c.carbonService.GetVehicleLogs(ctx.Context(), userID, vehicleID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "vehicle logs retrieved successfully", logs))
}

func (c *CarbonController) CreateElectronic(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	dto := new(dto.CreateElectronicDTO)
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid body"))
	}

	electronic, err := c.carbonService.CreateElectronic(ctx.Context(), userID, dto)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(helpers.SuccessResponseWithData(true, "electronic device created successfully", electronic))
}

func (c *CarbonController) ListUserElectronics(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	electronics, err := c.carbonService.ListUserElectronics(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "electronic devices retrieved successfully", electronics))
}

func (c *CarbonController) AddElectronicsLog(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	dto := new(dto.AddElectronicsLogDTO)
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid body"))
	}

	if err := c.carbonService.AddElectronicsLog(ctx.Context(), userID, dto); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.BasicResponse(true, "electronics log berhasil ditambahkan"))
}

func (c *CarbonController) GetElectronicsLogs(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	deviceID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid device ID"))
	}

	logs, err := c.carbonService.GetElectronicsLogs(ctx.Context(), userID, deviceID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "electronics logs retrieved successfully", logs))
}


func (c *CarbonController) EditVehicle(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	vehicleID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid vehicle ID"))
	}

	dto := new(dto.EditVehicleDTO)
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid body"))
	}

	vehicle, err := c.carbonService.EditVehicle(ctx.Context(), userID, vehicleID, dto)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "vehicle updated successfully", vehicle))
}

func (c *CarbonController) DeleteVehicle(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	vehicleID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid vehicle ID"))
	}

	if err := c.carbonService.DeleteVehicle(ctx.Context(), userID, vehicleID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.BasicResponse(true, "vehicle deleted successfully"))	
}

func (c *CarbonController) EditElectronic(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	deviceID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid device ID"))
	}

	dto := new(dto.EditElectronicDTO)
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid body"))
	}

	electronic, err := c.carbonService.EditElectronic(ctx.Context(), userID, deviceID, dto)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "electronic device updated successfully", electronic))
}

func (c *CarbonController) DeleteElectronic(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	deviceID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "invalid device ID"))
	}

	if err := c.carbonService.DeleteElectronic(ctx.Context(), userID, deviceID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.BasicResponse(true, "electronic device deleted successfully"))
}

func (c *CarbonController) GetAllVehicleLogs(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	logs, err := c.carbonService.GetAllVehicleLogs(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "all vehicle logs retrieved successfully", logs))
}

func (c *CarbonController) GetAllElectronicLogs(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)

	logs, err := c.carbonService.GetAllElectronicLogs(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "all electronic logs retrieved successfully", logs))
}
