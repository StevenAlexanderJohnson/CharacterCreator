package main

import (
	"dndcc/internal"
	"dndcc/internal/controllers"
	"dndcc/internal/database"
	"dndcc/internal/models"
	"dndcc/internal/repositories"
	"dndcc/internal/services"
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
	_ "modernc.org/sqlite"
)

func main() {
	app := grove.NewApp("ccapi")

	config, err := internal.ParseAppConfig()
	if err != nil {
		panic(err)
	}

	logger := grove.NewDefaultLogger("ccapi-auth")
	authConfig, err := grove.LoadAuthenticatorConfigFromEnv()
	if err != nil {
		panic(err)
	}

	db, err := database.CreateDatabaseConnection(config.DB.SqlitePath)
	if err != nil {
		panic(err)
	}

	authenticator := grove.NewAuthenticator[*models.Claims](authConfig)

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, authenticator)

	sessionRepo := repositories.NewSessionRepository(db)
	sessionService := services.NewSessionService(sessionRepo)

	app.
		WithMiddleware(grove.DefaultRequestLoggerMiddleware(logger)).
		WithController(controllers.NewAuthController(authService, sessionService, logger)).
		WithController(controllers.NewHomeController(logger, authenticator)).
		WithRoute("/public/", http.FileServer(http.Dir("public")))

	if err := app.Run(); err != nil {
		panic(err)
	}
}
