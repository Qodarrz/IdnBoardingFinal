// routes/setup.go
package routes

import (
	"database/sql"

	"github.com/Qodarrz/fiber-app/controller"
	"github.com/Qodarrz/fiber-app/middleware"
	"github.com/Qodarrz/fiber-app/repository"
	"github.com/Qodarrz/fiber-app/service"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, db *sql.DB, mw *middleware.Middlewares) {
	authService := service.NewAuthService(
		repository.NewUserRepository(db),
		repository.NewActivityRepository(db),
		repository.CheckMissionRepository(db),
	)

	carbonService := service.NewCarbonService(
		repository.NewCarbonRepository(db),
		repository.CheckMissionRepository(db),
	)

	missionRepo := repository.NewMissionRepository(db)
	userMissionRepo := repository.NewMissionRepository(db)

	userMissionService := service.NewMissionService(
		missionRepo,
		userMissionRepo,
		repository.NewBadgeRepository(db),
	)

	storeRepo := repository.NewStoreRepository(db)
	pointsRepo := repository.NewPointsRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	notificationRepo := repository.NewNotificationRepo(db)

	storeService := service.NewStoreService(
		storeRepo,
		pointsRepo,
		activityRepo,
		notificationRepo,
	)

	userCustomService := service.NewUserCustomEndpointService(
		repository.NewUserCustomEndpointRepo(db),
		repository.CheckMissionRepository(db),
		repository.NewMissionRepository(db),
	)

	profileService := service.NewUserProfileService(
		repository.NewUserProfileRepository(db),
		repository.NewActivityRepository(db),
		repository.NewUserRepository(db),
		
	)

	chatbotService, err := service.NewGeminiService()
	if err != nil {
		println("⚠️ Gagal init GeminiService:", err.Error())

		api := app.Group("/api/gemini")
		api.Post("/generate", func(c *fiber.Ctx) error {
			return c.Status(500).JSON(fiber.Map{
				"error": "GeminiService not initialized: " + err.Error(),
			})
		})
	} else {
		controller.InitGeminiController(app, chatbotService, mw)
	}

	badgeService := service.NewBadgeService(repository.NewBadgeRepository(db))
	notifCustomService := service.NewNotificationService(repository.NewNotificationRepo(db))

	controller.InitAuthController(app, authService, mw)
	controller.InitCarbonController(app, carbonService, mw)
	controller.InitMissionController(app, userMissionService, mw)
	controller.InitStoreController(app, storeService, mw)
	controller.InitBadgeController(app, badgeService, mw)
	controller.InitUserProfileController(app, profileService, mw)
	controller.InitUserCustomEndpointController(app, userCustomService, notifCustomService, mw)

}
