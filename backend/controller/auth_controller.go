package controller

import (
	"fmt"
	"net/http"
	"strconv"

	dto "github.com/Qodarrz/fiber-app/dto"
	helpers "github.com/Qodarrz/fiber-app/helper"
	"github.com/Qodarrz/fiber-app/middleware"
	service "github.com/Qodarrz/fiber-app/service"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService service.AuthServiceInterface
}

func InitAuthController(app *fiber.App, svc service.AuthServiceInterface, mw *middleware.Middlewares) {
	ctrl := &AuthController{authService: svc}

	public := app.Group("/api/auth")
	public.Post("/register", ctrl.Register)
	public.Post("/login", ctrl.Login)
	public.Post("/google", ctrl.GoogleLogin)
	public.Get("/verify-email", ctrl.VerifyEmail)
	public.Post("/verify-email", ctrl.VerifyEmail)
	public.Post("/reset-password/request", ctrl.RequestResetPassword)
	public.Post("/reset-password", ctrl.ResetPassword)

	// private := app.Group("/api/auth", mw.JWT, middleware.VerifyEmailMiddleware(mw.DB))
	private := app.Group("/api/auth", mw.JWT)
	private.Get("/me", ctrl.Profile)
	private.Post("/update-password", ctrl.UpdatePassword)
	private.Post("/link/google", ctrl.LinkGoogleAccount)
	private.Post("/unlink/google", ctrl.UnlinkGoogleAccount)
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	req := new(dto.RegisterDTO)
	if err := helpers.BindAndValidate(ctx, req); err != nil {
		if vErr, ok := err.(*helpers.ValidationError); ok {
			return ctx.Status(http.StatusBadRequest).JSON(helpers.ErrorResponseRequest(false, vErr.Message, vErr.Errors))
		}
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	user, err := c.authService.Register(ctx.Context(), req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	resp := dto.RegisterResponseDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "registrasi berhasil", resp))
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	req := new(dto.LoginDTO)
	if err := helpers.BindAndValidate(ctx, req); err != nil {
		if vErr, ok := err.(*helpers.ValidationError); ok {
			return ctx.Status(http.StatusBadRequest).JSON(helpers.ErrorResponseRequest(false, vErr.Message, vErr.Errors))
		}
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	user, token, err := c.authService.Login(ctx.Context(), req)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, err.Error()))
	}

	resp := dto.AuthResponseDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "login berhasil", resp))
}

func (c *AuthController) Profile(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(helpers.BasicResponse(false, "token tidak valid"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(helpers.BasicResponse(false, "user ID tidak valid"))
	}

	user, err := c.authService.GetProfile(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(helpers.SuccessResponseWithData(true, "profil ditemukan", user))
}

func (c *AuthController) VerifyEmail(ctx *fiber.Ctx) error {
	fmt.Println("FULL URL:", ctx.OriginalURL())
    fmt.Println("QUERY token:", ctx.Query("token"))

    token := ctx.Query("token")
    if token == "" {
        return ctx.Status(http.StatusBadRequest).JSON(
            helpers.ErrorResponseRequest(false, "Bad Request", map[string]string{
                "token": "token wajib diisi",
            }),
        )
    }
    if token == "" {
        var body struct {
            Token string `json:"token"`
        }
        if err := ctx.BodyParser(&body); err != nil {
            return ctx.Status(http.StatusBadRequest).JSON(
                helpers.BasicResponse(false, "token wajib diisi"),
            )
        }
        token = body.Token
    }

    if token == "" {
        return ctx.Status(http.StatusBadRequest).JSON(
            helpers.ErrorResponseRequest(false, "Bad Request", map[string]string{
                "token": "token wajib diisi",
            }),
        )
    }

    userID, err := helpers.DecodeJWT(token)
    if err != nil {
        return ctx.Status(http.StatusBadRequest).JSON(
            helpers.BasicResponse(false, err.Error()),
        )
    }

    if err := c.authService.VerifyEmail(ctx.Context(), userID); err != nil {
        return ctx.Status(http.StatusBadRequest).JSON(
            helpers.BasicResponse(false, err.Error()),
        )
    }

    return ctx.Status(http.StatusOK).JSON(
        helpers.SuccessResponseWithData(true, "email berhasil diverifikasi", nil),
    )
}

func (c *AuthController) RequestResetPassword(ctx *fiber.Ctx) error {
	req := new(dto.ResetPasswordRequestDTO)
	if err := helpers.BindAndValidate(ctx, req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	if err := c.authService.RequestResetPassword(ctx.Context(), req.Email); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "link reset password terkirim", nil))
}

func (c *AuthController) ResetPassword(ctx *fiber.Ctx) error {
	req := new(dto.ResetPasswordDTO)
	if err := helpers.BindAndValidate(ctx, req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	if err := c.authService.ResetPassword(ctx.Context(), req.Token, req.NewPassword); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "password berhasil diubah", nil))
}

func (c *AuthController) UpdatePassword(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "token tidak valid"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "user ID tidak valid"))
	}

	req := new(dto.UpdatePasswordDTO)
	if err := helpers.BindAndValidate(ctx, req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	if err := c.authService.UpdatePassword(ctx.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "password berhasil diperbarui", nil))
}

func (c *AuthController) GoogleLogin(ctx *fiber.Ctx) error {
	req := new(dto.GoogleAuthRequest)
	if err := helpers.BindAndValidate(ctx, req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	// Verifikasi token Google
	googleUser, err := helpers.VerifyGoogleToken(req.Token)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "Token Google tidak valid"))
	}

	
		googleUserDTO := &dto.GoogleUserDTO{
		ID:      googleUser.ID,
		Email:   googleUser.Email,
		Name:    googleUser.Name,
		Picture: googleUser.Picture,
	}

	user, token, err := c.authService.LoginWithGoogle(ctx.Context(), googleUserDTO)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(helpers.BasicResponse(false, err.Error()))
	}

	resp := dto.AuthResponseDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "login dengan Google berhasil", resp))
}
func (c *AuthController) LinkGoogleAccount(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "token tidak valid"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "user ID tidak valid"))
	}

	req := new(dto.GoogleAuthRequest)
	if err := helpers.BindAndValidate(ctx, req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	// Verifikasi token Google
	googleUser, err := helpers.VerifyGoogleToken(req.Token)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "Token Google tidak valid"))
	}

	// Konversi ke GoogleUserDTO
	googleUserDTO := &dto.GoogleUserDTO{
		ID:      googleUser.ID,
		Email:   googleUser.Email,
		Name:    googleUser.Name,
		Picture: googleUser.Picture,
	}

	if err := c.authService.LinkGoogleAccount(ctx.Context(), userID, googleUserDTO); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "akun Google berhasil dihubungkan", nil))
}

func (c *AuthController) UnlinkGoogleAccount(ctx *fiber.Ctx) error {
	claims := helpers.GetUserClaims(ctx)
	if claims == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "token tidak valid"))
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, "user ID tidak valid"))
	}

	if err := c.authService.UnlinkGoogleAccount(ctx.Context(), userID); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(helpers.BasicResponse(false, err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(helpers.SuccessResponseWithData(true, "akun Google berhasil diputus", nil))
}