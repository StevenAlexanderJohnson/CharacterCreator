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

func setAuthCookie(w http.ResponseWriter, jwt string) {
	cookie := &http.Cookie{
		Name:     string(authCookie),
		Value:    jwt,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

func getAuthCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(string(authCookie))
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
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

func clearAuthCookies(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     string(sessionCookie),
		Value:    "",
		Path:     "/",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	cookie = &http.Cookie{
		Name:     string(authCookie),
		Value:    "",
		Path:     "/",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}
