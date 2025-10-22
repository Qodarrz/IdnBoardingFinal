// controller/user_profile_controller.go
package controller

import (
	"fmt"
	"log"
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

type UserProfileController struct {
	userProfileService service.UserProfileServiceInterface
}

func InitUserProfileController(app *fiber.App, svc service.UserProfileServiceInterface, mw *middleware.Middlewares) {
	ctrl := &UserProfileController{userProfileService: svc}

	private := app.Group("/api/user", mw.JWT)
	private.Get("/profile", ctrl.GetProfile)
	private.Patch("/profile", ctrl.UpdateProfile)
}

func (c *UserProfileController) GetProfile(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(helpers.BasicResponse(false, "token tidak valid"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "user ID tidak valid"))
	}

	profile, err := c.userProfileService.GetProfile(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "profil ditemukan", profile))
}

func (c *UserProfileController) UpdateProfile(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(helpers.BasicResponse(false, "token tidak valid"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "user ID tidak valid"))
	}

	req := &dto.UserProfileUpdateDTO{}

	// Parse manual semua field dari form-data (TANPA username)
	if v := ctx.FormValue("full_name"); v != "" {
		req.FullName = &v
	}
	if v := ctx.FormValue("gender"); v != "" {
		req.Gender = &v
	}
	if v := ctx.FormValue("birthdate"); v != "" {
		t, parseErr := time.Parse("2006-01-02", v)
		if parseErr != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(
				helpers.BasicResponse(false, "format birthdate tidak valid, gunakan YYYY-MM-DD"),
			)
		}
		req.Birthdate = &t
	}

	// Handle file upload avatar
	fileHeader, err := ctx.FormFile("avatar")
	if err == nil && fileHeader != nil {
		if fileHeader.Size > 5*1024*1024 {
			return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "ukuran file terlalu besar"))
		}

		contentType := fileHeader.Header.Get("Content-Type")
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/gif":  true,
			"image/jpg":  true,
		}
		
		if !allowedTypes[contentType] {
			return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "tipe file tidak didukung"))
		}

		// Generate unique filename
		fileExt := filepath.Ext(fileHeader.Filename)
		uniqueFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), "avatar", fileExt)
		tempPath := filepath.Join(os.TempDir(), uniqueFilename)

		if err := ctx.SaveFile(fileHeader, tempPath); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, "gagal simpan file sementara"))
		}
		defer os.Remove(tempPath)

		// Upload to storage
		url, err := helpers.UploadFile(tempPath)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, "gagal upload avatar"))
		}

		req.AvatarURL = &url
	}

	// Validasi minimal ada satu field yang diupdate (TANPA username)
	if req.FullName == nil && req.Gender == nil && req.Birthdate == nil && req.AvatarURL == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "tidak ada data yang diupdate"))
	}

	updatedProfile, err := c.userProfileService.UpdateProfile(ctx.Context(), userID, req)
	if err != nil {
		log.Printf("Error updating profile: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(helpers.BasicResponse(false, "gagal mengupdate profil"))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "profil berhasil diperbarui", updatedProfile))
}