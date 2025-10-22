// Di helper/auth_helper.go
package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Qodarrz/fiber-app/dto"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOAuthConfig *oauth2.Config
)

func InitGoogleOAuth(clientID, clientSecret, redirectURL string) {
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func VerifyGoogleToken(token string) (*dto.GoogleUserDTO, error) {
	// Verifikasi menggunakan Google API
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("gagal memverifikasi token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token tidak valid, status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response: %v", err)
	}

	var tokenInfo struct {
		Aud           string `json:"aud"`
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Locale        string `json:"locale"`
	}

	err = json.Unmarshal(body, &tokenInfo)
	if err != nil {
		return nil, fmt.Errorf("gagal parsing response: %v", err)
	}

	// Validasi client ID
	if tokenInfo.Aud != GoogleOAuthConfig.ClientID {
		return nil, fmt.Errorf("token tidak valid untuk aplikasi ini")
	}

	googleUser := &dto.GoogleUserDTO{
		ID:            tokenInfo.Sub,
		Email:         tokenInfo.Email,
		VerifiedEmail: tokenInfo.EmailVerified == "true",
		Name:          tokenInfo.Name,
		GivenName:     tokenInfo.GivenName,
		FamilyName:    tokenInfo.FamilyName,
		Picture:       tokenInfo.Picture,
		Locale:        tokenInfo.Locale,
	}

	return googleUser, nil
}