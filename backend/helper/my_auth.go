package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GetUserClaims(c *fiber.Ctx) *JWTClaims {
	u := c.Locals("user")
	if u == nil {
		return nil
	}

	token, ok := u.(*jwt.Token)
	if !ok || token == nil {
		return nil
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil
	}

	return claims
}

func GenerateEmailVerificationToken(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_EMAIL_SECRET")))
}

func GenerateJWT(userID string) (string, error) {
	now := time.Now()

	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	parts := strings.Split(signedToken, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("generated token has %d parts, expected 3", len(parts))
	}

	return signedToken, nil
}

// GenerateRandomToken menghasilkan token acak dengan panjang n byte
func GenerateRandomToken(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func DecodeJWT(tokenStr string) (int64, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_EMAIL_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("token tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("tidak bisa membaca claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return 0, errors.New("user_id tidak ada di token")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, errors.New("user_id token invalid")
	}

	return userID, nil
}


func UploadFile(filePath string) (string, error) {
    // buka file lokal
    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("gagal buka file: %w", err)
    }
    defer file.Close()

    // buat buffer untuk multipart form
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    // tambahkan file ke form
    part, err := writer.CreateFormFile("file", file.Name())
    if err != nil {
        return "", fmt.Errorf("gagal create form file: %w", err)
    }
    if _, err := io.Copy(part, file); err != nil {
        return "", fmt.Errorf("gagal copy file: %w", err)
    }
    writer.Close()

    // request POST ke endpoint Laravel
    req, err := http.NewRequest("POST", "https://sipdesaraksajaya.my.id/api/storage", body)
    if err != nil {
        return "", fmt.Errorf("gagal bikin request: %w", err)
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("gagal kirim request: %w", err)
    }
    defer resp.Body.Close()

    // baca response body
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("gagal baca response: %w", err)
    }

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return "", fmt.Errorf("upload gagal: %s", string(respBody))
    }

    // parsing response JSON untuk dapat URL
    type respJSON struct {
        Success bool   `json:"success"`
        Url     string `json:"url"`
    }
    var r respJSON
    if err := json.Unmarshal(respBody, &r); err != nil {
        return "", fmt.Errorf("gagal parse JSON: %w", err)
    }

    return r.Url, nil
}