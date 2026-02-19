package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"
	"strings"
)

const (
	SessionCookieName = "playtime_session"
	SessionSecretEnv  = "SESSION_SECRET"
)

func getSecret() []byte {
	s := os.Getenv(SessionSecretEnv)
	if s == "" {
		s = "dev-secret-change-in-production"
	}
	return []byte(s)
}

func sign(value string) string {
	mac := hmac.New(sha256.New, getSecret())
	mac.Write([]byte(value))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func SetSteamID(w http.ResponseWriter, steamID string) {
	encoded := base64.URLEncoding.EncodeToString([]byte(steamID))
	sig := sign(steamID)
	value := encoded + "." + sig
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30, // 30 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func GetSteamID(r *http.Request) string {
	c, err := r.Cookie(SessionCookieName)
	if err != nil || c.Value == "" {
		return ""
	}
	parts := strings.SplitN(c.Value, ".", 2)
	if len(parts) != 2 {
		return ""
	}
	decoded, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return ""
	}
	steamID := string(decoded)
	expectedSig, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return ""
	}
	mac := hmac.New(sha256.New, getSecret())
	mac.Write([]byte(steamID))
	if hmac.Equal(mac.Sum(nil), expectedSig) {
		return steamID
	}
	return ""
}

func ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// BaseURL returns the app base URL for redirects (e.g. https://localhost:8080).
func BaseURL(r *http.Request) string {
	if u := os.Getenv("BASE_URL"); u != "" {
		return strings.TrimSuffix(u, "/")
	}
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	return scheme + "://" + r.Host
}
