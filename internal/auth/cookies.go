package auth

import (
	"net/http"
	"time"
)


func setCookieHandler(w http.ResponseWriter, r *http.Request, refresh string) {
    cookies := &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false,  // only for local development
		SameSite: http.SameSiteStrictMode,
		Path: "/",
	}

    http.SetCookie(w, cookies)
}


func clearCookieHandler(w http.ResponseWriter, r *http.Request) { 
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0), 
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}