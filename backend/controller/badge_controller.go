package controller

import (
	"net/http"
	"strconv"

	helper "github.com/Qodarrz/fiber-app/helper"
	"github.com/Qodarrz/fiber-app/middleware"
	"github.com/Qodarrz/fiber-app/service"
	"github.com/gofiber/fiber/v2"
)

type BadgeController struct {
	badgeService service.BadgeService
}

func InitBadgeController(app *fiber.App, svc service.BadgeService, mw *middleware.Middlewares) {
	ctrl := &BadgeController{badgeService: svc}

	public := app.Group("/api/badges")
	private := public.Group("/", mw.JWT)

	private.Get("/", ctrl.GetUserBadges)
}

func (c *BadgeController) GetUserBadges(ctx *fiber.Ctx) error {
	claims := helper.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(helper.BasicResponse(false, "token tidak valid"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helper.BasicResponse(false, "user ID tidak valid"))
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	badges, err := c.badgeService.GetBadgesWithOwnership(ctx.Context(), userID, page, limit)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helper.BasicResponse(false, "gagal mengambil badge"))
	}

	return ctx.Status(http.StatusOK).JSON(helper.SuccessResponseWithData(true, "badge berhasil diambil", fiber.Map{
		"data":  badges,
		"page":  page,
		"limit": limit,
	}))
}
