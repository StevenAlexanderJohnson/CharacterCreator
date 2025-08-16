package controllers

import (
	"errors"
	"net/http"
	"time"
)

type cookieName string

var (
	sessionCookie cookieName = "session"
)

func setAuthCookie(jwt string) {

}

func getSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(string(sessionCookie))
	if err != nil {
		if errors.Is(http.ErrNoCookie, err) {
			return "", nil
		}
		return "", err
	}
	return cookie.Value, nil
}

func setSessionCookie(w http.ResponseWriter, session string) {
	cookie := &http.Cookie{
		Name:     string(sessionCookie),
		Value:    session,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}
