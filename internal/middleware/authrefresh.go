package middleware

import (
	"context"
	"dndcc/internal/controllers"
	"dndcc/internal/models"
	"dndcc/internal/services"
	"net/http"
	"slices"
	"strings"

	"github.com/StevenAlexanderJohnson/grove"
)

type AuthWithRefreshMiddleware struct {
	logger          grove.ILogger
	authenticator   grove.Authenticator[*models.Claims]
	sessionService  *services.SessionService
	authService     *services.AuthService
	routeExceptions []string
}

func NewAuthWithRefreshMiddleware(
	logger grove.ILogger,
	authenticator grove.Authenticator[*models.Claims],
	sessionService *services.SessionService,
	authService *services.AuthService,
) *AuthWithRefreshMiddleware {
	return &AuthWithRefreshMiddleware{
		logger,
		authenticator,
		sessionService,
		authService,
		[]string{},
	}
}

func (a *AuthWithRefreshMiddleware) WithRouteException(route string) *AuthWithRefreshMiddleware {
	a.routeExceptions = append(a.routeExceptions, route)
	return a
}

func (a *AuthWithRefreshMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public/") {
			next.ServeHTTP(w, r)
			return
		}
		authCookie, err := r.Cookie("session_token")
		if err == nil && authCookie.Value != "" {
			claims := &models.Claims{}
			_, err := a.authenticator.VerifyToken(authCookie.Value, claims)
			if err == nil {
				authContext := context.WithValue(r.Context(), grove.AuthTokenKey, claims)
				next.ServeHTTP(w, r.WithContext(authContext))
				return
			}
		}

		refreshTokenCookie, err := r.Cookie("session")
		if err != nil || refreshTokenCookie.Value == "" {
			if slices.Contains(a.routeExceptions, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		session, err := a.sessionService.Get(refreshTokenCookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authToken, lifetime, err := a.authService.GetTokenById(session.UserId)
		if err != nil {
			a.logger.Error("an error occurred while automatically refreshing user token: %v", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		controllers.SetAuthCookie(w, authToken, lifetime)

		claims := &models.Claims{}
		parsedClaims, _ := a.authenticator.VerifyToken(authToken, claims)

		authContext := context.WithValue(r.Context(), grove.AuthTokenKey, parsedClaims)
		next.ServeHTTP(w, r.WithContext(authContext))
	})
}
