package main

import (
	"context"
	"dndcc/internal"
	"dndcc/internal/controllers"
	"dndcc/internal/database"
	"dndcc/internal/middleware"
	"dndcc/internal/models"
	"dndcc/internal/repositories"
	"dndcc/internal/services"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/StevenAlexanderJohnson/grove"
	_ "modernc.org/sqlite"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
	defer db.Close()

	authenticator := grove.NewAuthenticator[*models.Claims](authConfig)

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo, authenticator)

	sessionRepo := repositories.NewSessionRepository(db)
	sessionService := services.NewSessionService(sessionRepo)

	characterRepo := repositories.NewCharacterRepository(db)
	characterService := services.NewCharacterService(characterRepo)

	authWithRefreshMiddleware := middleware.
		NewAuthWithRefreshMiddleware(logger, *authenticator, sessionService, authService).
		WithRouteException("/").WithRouteException("/auth/login").WithRouteException("/auth/register")

	authScope := grove.NewScope().
		WithMiddleware(authWithRefreshMiddleware.Middleware).
		WithController(controllers.NewAuthController(authService, sessionService, logger)).
		WithController(controllers.NewHomeController(logger, authenticator)).
		WithController(controllers.NewCharacterController(logger, characterService))
	app.
		WithScope("/", authScope).
		WithRoute("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})).
		WithRoute("/public/", http.FileServer(http.Dir("public")))

	go func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down gracefully, press Ctrl+C again to force")
	db.Close()
	logger.Info("database connection closed")
}
