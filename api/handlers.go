package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	moviehandlers "github.com/samthehai/ml-backend-test-samthehai/internal/movie/interfaceadapters/http"
	movierepository "github.com/samthehai/ml-backend-test-samthehai/internal/movie/interfaceadapters/repository"
	movieusecase "github.com/samthehai/ml-backend-test-samthehai/internal/movie/usecase"
	userhandlers "github.com/samthehai/ml-backend-test-samthehai/internal/user/interfaceadapters/http"
	userrepository "github.com/samthehai/ml-backend-test-samthehai/internal/user/interfaceadapters/repository"
	userusecase "github.com/samthehai/ml-backend-test-samthehai/internal/user/usecase"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/csrf"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/token"
	"github.com/samthehai/ml-backend-test-samthehai/pkg/utils"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// repository
	userRepository := userrepository.NewUserRepository(s.connManager)
	movieRepository := movierepository.NewMovieRepository(s.connManager)

	tokenMaker, err := token.NewJWTMaker(s.cfg.Server.JWTSecretKey)
	if err != nil {
		return err
	}

	// usecase
	userUsecase := userusecase.NewUserUsecase(*s.cfg, userRepository, s.logger, tokenMaker)
	movieUsecase := movieusecase.NewMovieUsecase(*s.cfg, movieRepository, s.logger)

	// handlers
	userHanlders := userhandlers.NewUserHandlers(s.cfg, userUsecase, s.logger)
	movieHanlders := moviehandlers.NewMovieHandlers(s.cfg, movieUsecase, s.logger)

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
	movieGroup.GET("/:id", movieHanlders.GetByID())

	// health check api
	health := v1.Group("/health")
	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
