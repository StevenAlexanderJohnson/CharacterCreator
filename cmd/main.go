package main

import (
	"dndcc/internal/controllers"
	"dndcc/internal/database"
	"dndcc/internal/models"
	"dndcc/internal/repositories"
	"dndcc/internal/services"

	"github.com/StevenAlexanderJohnson/grove"
	_ "modernc.org/sqlite"
)

func main() {
	app := grove.NewApp("ccapi")

	logger := grove.NewDefaultLogger("ccapi-auth")
	authConfig, err := grove.LoadAuthenticatorConfigFromEnv()
	if err != nil {
		panic(err)
	}

	db, err := database.CreateDatabaseConnection("database.db")
	if err != nil {
		panic(err)
	}

	authenticator := grove.NewAuthenticator[models.Claims](authConfig)

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, authenticator)

	authScope := grove.NewScope().
		WithMiddleware(
			grove.DefaultAuthMiddleware(
				authenticator,
				logger,
				func() models.Claims {
					return models.Claims{}
				},
			),
		)
	app.
		WithScope("auth", authScope).
		WithController(controllers.NewAuthController(authService, logger))

	if err := app.Run(); err != nil {
		panic(err)
	}
}
