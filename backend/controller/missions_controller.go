// controller/mission.go
package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	dto "github.com/Qodarrz/fiber-app/dto"
	helpers "github.com/Qodarrz/fiber-app/helper"
	"github.com/Qodarrz/fiber-app/middleware"
	service "github.com/Qodarrz/fiber-app/service"
	"github.com/gofiber/fiber/v2"
)

type MissionController struct {
	missionService service.MissionServiceInterface
}

func InitMissionController(app *fiber.App, svc service.MissionServiceInterface, mw *middleware.Middlewares) {
	ctrl := &MissionController{missionService: svc}

	// Public routes
	public := app.Group("/api/missions")
	public.Get("/", ctrl.GetAllMissions)
	public.Get("/active", ctrl.GetActiveMissions)
	public.Get("/:id", ctrl.GetMissionByID)

	// Private routes (require authentication)
	private := app.Group("/api/missions", mw.JWT)
	private.Post("/", ctrl.CreateMission)
	private.Get("/my-missions", ctrl.GetUserMissions)
	private.Post("/with-badge", ctrl.CreateMissionWithBadge)
	private.Get("/:id/check-completion", ctrl.CheckMissionCompletion)
}

func (c *MissionController) CreateMission(ctx *fiber.Ctx) error {
	req := new(dto.CreateMissionDTO)
	if err := helpers.BindAndValidate(ctx, req); err != nil {
		if vErr, ok := err.(*helpers.ValidationError); ok {
			return ctx.Status(http.StatusBadRequest).JSON(helpers.ErrorResponseRequest(false, vErr.Message, vErr.Errors))
		}
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	mission, err := c.missionService.CreateMission(ctx.Context(), req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusCreated).JSON(helpers.SuccessResponseWithData(true, "Mission created successfully", mission))
}

func (c *MissionController) GetMissionByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "Invalid mission ID"))
	}

	mission, err := c.missionService.GetMissionByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "Mission found", mission))
}

func (c *MissionController) GetAllMissions(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	missions, err := c.missionService.GetAllMissions(ctx.Context(), page, limit)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "Missions retrieved successfully", missions))
}

func (c *MissionController) GetActiveMissions(ctx *fiber.Ctx) error {
	missions, err := c.missionService.GetActiveMissions(ctx.Context())
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "Active missions retrieved successfully", missions))
}

func (c *MissionController) GetUserMissions(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "Invalid token"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "Invalid user ID"))
	}

	missions, err := c.missionService.GetUserMissions(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "User missions retrieved successfully", missions))
}

func (c *MissionController) CreateMissionWithBadge(ctx *fiber.Ctx) error {
    req := new(dto.CreateMissionWithBadgeDTO)

    // Ambil field manual dari form-data
    req.Title = ctx.FormValue("title")
    req.Description = ctx.FormValue("description")
    req.MissionType = dto.MissionType(ctx.FormValue("mission_type"))
    req.CriteriaType = dto.CriteriaType(ctx.FormValue("criteria_type"))
    req.PointsReward, _ = strconv.Atoi(ctx.FormValue("points_reward"))
    req.GivesBadge, _ = strconv.ParseBool(ctx.FormValue("gives_badge"))
    req.BadgeName = ctx.FormValue("badge_name")
    req.BadgeDescription = ctx.FormValue("badge_description")
    req.TargetValue, _ = strconv.ParseFloat(ctx.FormValue("target_value"), 64)
    if exp := ctx.FormValue("expired_at"); exp != "" {
        t, _ := time.Parse(time.RFC3339, exp)
        req.ExpiredAt = &t
    }

    // Handle file upload
    fileHeader, err := ctx.FormFile("file")
if err == nil && fileHeader != nil {
    tempPath := filepath.Join(os.TempDir(), fileHeader.Filename)
    if err := ctx.SaveFile(fileHeader, tempPath); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
    }
    url, err := helpers.UploadFile(tempPath)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
    }
    req.BadgeImageURL = url
}


    // Validasi manual jika perlu
    if req.Title == "" || req.MissionType == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "title atau mission_type wajib diisi"})
    }

    result, err := c.missionService.CreateMissionWithBadge(ctx.Context(), req)
    if err != nil {
        return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
    }

    return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Mission with badge created successfully", "data": result})
}


func (c *MissionController) CheckMissionCompletion(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "Invalid token"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "Invalid user ID"))
	}

	missionID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "Invalid mission ID"))
	}

	completed, err := c.missionService.CheckMissionCompletion(ctx.Context(), userID, missionID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "Mission completion status", map[string]bool{
		"completed": completed,
	}))
}
