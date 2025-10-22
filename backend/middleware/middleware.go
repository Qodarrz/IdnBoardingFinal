package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	helpers "github.com/Qodarrz/fiber-app/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

type Middlewares struct {
	KeyApi fiber.Handler
	JWT    fiber.Handler
	DB     *sql.DB
}

func InitRateLimiterConfig() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        30,
		Expiration: 3 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(http.StatusTooManyRequests).JSON(
				helpers.BasicResponse(false, "terlalu banyak request"),
			)
		},
	})
}

func AuthKeyMiddleware(apiKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("X-API-Key")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "X-API-Key header missing"))
		}

		if authHeader != apiKey {
			return c.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, "Invalid API Key"))
		}

		return c.Next()
	}
}

func InitMiddlewares(db *sql.DB) *Middlewares {
	return &Middlewares{
		JWT: initJWTMiddleware(),
		DB:  db,
	}
}

func initJWTMiddleware() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable not set")
	}

	return jwtware.New(jwtware.Config{
		SigningKey: []byte(secret),
		Claims:     &helpers.JWTClaims{},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			errorMsg := "Token error"
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				errorMsg = "Token format invalid"
			case errors.Is(err, jwt.ErrTokenExpired):
				errorMsg = "Token expired"
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				errorMsg = "Token not yet valid"
			default:
				errorMsg = fmt.Sprintf("Token error: %v", err)
			}
			return c.Status(http.StatusUnauthorized).JSON(helpers.BasicResponse(false, errorMsg))
		},
	})
}

func VerifyEmailMiddleware(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := helpers.GetUserClaims(c)
		if claims == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helpers.BasicResponse(false, "invalid token claims"),
			)
		}

		userID, err := strconv.ParseInt(claims.UserID, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helpers.BasicResponse(false, "invalid user id in token"),
			)
		}

		var emailVerifiedAt sql.NullTime
		err = db.QueryRow("SELECT email_verified_at FROM users WHERE id = $1", userID).Scan(&emailVerifiedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(fiber.StatusUnauthorized).JSON(
					helpers.BasicResponse(false, "user tidak ditemukan"),
				)
			}
			fmt.Println("DB error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				helpers.BasicResponse(false, "internal server error"),
			)
		}

		if !emailVerifiedAt.Valid {
			return c.Status(fiber.StatusForbidden).JSON(
				helpers.BasicResponse(false, "email belum diverifikasi"),
			)
		}

		return c.Next()
	}
}



func AdminMiddleware(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken := c.Locals("user")
		if userToken == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helpers.BasicResponse(false, "unauthorized"),
			)
		}

		claims, ok := userToken.(*helpers.JWTClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helpers.BasicResponse(false, "invalid token claims"),
			)
		}

		userIDStr := claims.ID
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helpers.BasicResponse(false, "invalid user id in token"),
			)
		}

		var role string
		err = db.QueryRow("SELECT role FROM users WHERE id = $1", userID).Scan(&role)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(fiber.StatusUnauthorized).JSON(
					helpers.BasicResponse(false, "user tidak ditemukan"),
				)
			}
			fmt.Println("DB error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				helpers.BasicResponse(false, "internal server error"),
			)
		}

		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(
				helpers.BasicResponse(false, "akses hanya untuk admin"),
			)
		}

		return c.Next()
	}
}
