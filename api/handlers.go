package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	docs "github.com/samthehai/ml-backend-test-samthehai/docs"
	"github.com/samthehai/ml-backend-test-samthehai/internal/middlewares"
	moviehandlers "github.com/samthehai/ml-backend-test-samthehai/internal/movie/interfaceadapters/http"
	movierepository "github.com/samthehai/ml-backend-test-samthehai/internal/movie/interfaceadapters/repository"
	movieusecase "github.com/samthehai/ml-backend-test-samthehai/internal/movie/usecase"
	userhandlers "github.com/samthehai/ml-backend-test-samthehai/internal/user/interfaceadapters/http"
	userrepository "github.com/samthehai/ml-backend-test-samthehai/internal/user/interfaceadapters/repository"
	userusecase "github.com/samthehai/ml-backend-test-samthehai/internal/user/usecase"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/csrf"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/logger"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/token"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/utils"
	echoswagger "github.com/swaggo/echo-swagger"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// repository
	userRepository := userrepository.NewUserRepository(s.connManager)
	movieRepository := movierepository.NewMovieRepository(s.connManager)
	favoriteRepository := movierepository.NewFavoriteRepository(s.connManager)

	tokenMaker, err := token.NewJWTMaker(s.cfg.Server.JWTSecretKey)
	if err != nil {
		return err
	}

	// usecase
	userUsecase := userusecase.NewUserUsecase(*s.cfg, userRepository, s.logger, tokenMaker)
	movieUsecase := movieusecase.NewMovieUsecase(*s.cfg, s.logger, movieRepository, favoriteRepository)

	// middlewares
	middlewareManager := middlewares.NewMiddlewareManager(s.cfg, s.logger, userUsecase)

	// handlers
	userHanlders := userhandlers.NewUserHandlers(s.cfg, userUsecase, s.logger)
	movieHanlders := moviehandlers.NewMovieHandlers(s.cfg, movieUsecase, s.logger, middlewareManager.GetCurrentUser)

	docs.SwaggerInfo.Title = "MonstarLab Backend Test REST API"
	docs.SwaggerInfo.BasePath = "/api/v1"
	e.GET("/swagger/*", echoswagger.WrapHandler)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	v1 := e.Group("/api/v1")

	// user api
	userGroup := v1.Group("/users")
	userGroup.POST("/register", userHanlders.Register())
	userGroup.POST("/login", userHanlders.Login())

	// movie api
	movieGroup := v1.Group("/movies")
	movieGroup.GET("", movieHanlders.SearchByKeyword())
	movieGroup.GET("/:id", movieHanlders.GetByID())

	// favorite api
	favoriteGroup := v1.Group("/favorites", middlewareManager.AuthMiddleware(tokenMaker))
	favoriteGroup.GET("", movieHanlders.ListFavoriteMovies())
	favoriteGroup.POST("/:id", movieHanlders.AddFavoriteMovie())

	// health check api
	health := v1.Group("/health")
	health.GET("", healthCheck(s.logger))

	return nil
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func healthCheck(logger logger.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	}
}
