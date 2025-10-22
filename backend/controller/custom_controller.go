package controller

import (
	"net/http"
	"strconv"

	dto "github.com/Qodarrz/fiber-app/dto"
	helper "github.com/Qodarrz/fiber-app/helper"
	"github.com/Qodarrz/fiber-app/middleware"
	"github.com/Qodarrz/fiber-app/service"
	"github.com/gofiber/fiber/v2"
)

type UserCustomEndpointController struct {
	userCustomService service.UserCustomEndpointServiceInterface
	notifcustomservice service.NotificationService
}

func InitUserCustomEndpointController(app *fiber.App, svc service.UserCustomEndpointServiceInterface, notifSvc service.NotificationService, mw *middleware.Middlewares) {
	ctrl := &UserCustomEndpointController{
		userCustomService: svc,
		notifcustomservice: notifSvc,
	}

	public := app.Group("/api/custom")
	public.Get("/leaderboard", ctrl.GetLeaderboard)

	private := app.Group("/api/custom", mw.JWT)
	private.Get("/user-data/:id", ctrl.GetUserCustomData)
	private.Get("/my-data", ctrl.GetMyCustomData)
	private.Get("/mission-progress", ctrl.GetMissionProgress)
	private.Get("/notifications", ctrl.GetNotifications)
}


func (c *UserCustomEndpointController) GetNotifications(ctx *fiber.Ctx) error {
	claims := helper.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helper.BasicResponse(false, "Unauthorized"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helper.BasicResponse(false, "Invalid user ID"))
	}

	notifications, err := c.notifcustomservice.GetNotificationsByUserID(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helper.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helper.SuccessResponseWithData(true, "Notifications retrieved successfully", notifications))
}


func (c *UserCustomEndpointController) GetMissionProgress(ctx *fiber.Ctx) error {
	claims := helper.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helper.BasicResponse(false, "Unauthorized"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helper.BasicResponse(false, "Invalid user ID"))
	}

	progressList, err := c.userCustomService.GetAllMissionProgress(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helper.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helper.SuccessResponseWithData(true, "Mission progress retrieved successfully", progressList))
}

func (c *UserCustomEndpointController) GetUserCustomData(ctx *fiber.Ctx) error {
	userID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helper.BasicResponse(false, "Invalid user ID"))
	}

	claims := helper.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helper.BasicResponse(false, "Unauthorized"))
	}

	currentUserID, _ := strconv.ParseInt(claims.UserID, 10, 64)
	if currentUserID != userID  {
		return ctx.Status(http.StatusForbidden).JSON(helper.BasicResponse(false, "Access denied"))
	}

	data, err := c.userCustomService.GetUserCustomData(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helper.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helper.SuccessResponseWithData(true, "User data retrieved successfully", data))
}

func (c *UserCustomEndpointController) GetMyCustomData(ctx *fiber.Ctx) error {
	claims := helper.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helper.BasicResponse(false, "Unauthorized"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helper.BasicResponse(false, "Invalid user ID"))
	}

	data, err := c.userCustomService.GetUserCustomData(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helper.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helper.SuccessResponseWithData(true, "User data retrieved successfully", data))
}

func (c *UserCustomEndpointController) GetLeaderboard(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)
	timeRange := ctx.Query("timeRange", "all")

	req := &dto.LeaderboardRequestDTO{
		Page:      page,
		Limit:     limit,
		TimeRange: timeRange,
	}

	leaderboard, err := c.userCustomService.GetLeaderboard(ctx.Context(), req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helper.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helper.SuccessResponseWithData(true, "Leaderboard retrieved successfully", leaderboard))
}
