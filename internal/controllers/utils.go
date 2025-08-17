package controllers

import (
	"errors"
	"net/http"
	"time"
)

type cookieName string

var (
	sessionCookie cookieName = "session"
	authCookie    cookieName = "session_token"
)

func SetAuthCookie(w http.ResponseWriter, jwt string, duration time.Duration) {
	cookie := &http.Cookie{
		Name:     string(authCookie),
		Value:    jwt,
		Path:     "/",
		Expires:  time.Now().Add(duration),
		MaxAge:   int(duration.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}

func getSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(string(sessionCookie))
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", nil
		}
		return "", err
	}
	return cookie.Value, nil
}

func SetSessionCookie(w http.ResponseWriter, session string, duration time.Duration) {
	cookie := &http.Cookie{
		Name:     string(sessionCookie),
		Value:    session,
		Path:     "/",
		Expires:  time.Now().Add(duration),
		MaxAge:   int(duration.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}

func clearAuthCookies(w http.ResponseWriter) {
	SetAuthCookie(w, "", -1)
	SetSessionCookie(w, "", -1)
}
