package auth

import (
	"net/http"
	"net/url"
	"regexp"
)

const steamOpenIDURL = "https://steamcommunity.com/openid/login"

var steamIDFromClaimedID = regexp.MustCompile(`https?://steamcommunity\.com/openid/id/(\d+)$`)

// LoginHandler redirects the user to Steam OpenID.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	base := BaseURL(r)
	returnTo := base + "/api/auth/steam/callback"
	realm := base
	params := url.Values{
		"openid.ns":         {"http://specs.openid.net/auth/2.0"},
		"openid.mode":      {"checkid_setup"},
		"openid.return_to": {returnTo},
		"openid.realm":     {realm},
		"openid.identity":  {"http://specs.openid.net/auth/2.0/identifier_select"},
		"openid.claimed_id": {"http://specs.openid.net/auth/2.0/identifier_select"},
	}
	redirectTo := steamOpenIDURL + "?" + params.Encode()
	http.Redirect(w, r, redirectTo, http.StatusFound)
}

// CallbackHandler handles the redirect from Steam: extract Steam ID from claimed_id, set session, redirect to /.
// We trust the redirect (Steam sent the user here with openid.claimed_id); no server-side check_authentication POST.
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Query().Get("openid.mode") != "id_res" {
		http.Redirect(w, r, BaseURL(r)+"/", http.StatusFound)
		return
	}
	claimedID := r.URL.Query().Get("openid.claimed_id")
	matches := steamIDFromClaimedID.FindStringSubmatch(claimedID)
	if len(matches) < 2 {
		http.Redirect(w, r, BaseURL(r)+"/", http.StatusFound)
		return
	}
	steamID := matches[1]
	SetSteamID(w, steamID)
	http.Redirect(w, r, BaseURL(r)+"/", http.StatusFound)
}

// MeHandler returns the current Steam ID or 401.
func MeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	steamID := GetSteamID(r)
	if steamID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"not logged in"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"steamId":"` + steamID + `"}`))
}

// LogoutHandler clears the session and redirects to /.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ClearSession(w)
	http.Redirect(w, r, BaseURL(r)+"/", http.StatusFound)
}
