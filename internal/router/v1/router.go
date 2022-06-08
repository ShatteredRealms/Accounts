package v1

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/internal/option"
	"github.com/ShatteredRealms/Accounts/pkg/repository"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"github.com/ShatteredRealms/GoUtils/pkg/middlewares"
	"github.com/ShatteredRealms/GoUtils/pkg/model"
	utilService "github.com/ShatteredRealms/GoUtils/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gorm.io/gorm"
)

// InitRouter Initializes all the routes for the service and starts the http server
func InitRouter(db *gorm.DB, config option.Config, logger log.LoggerService) (*gin.Engine, error) {
	if config.IsRelease() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(middlewares.ContentTypeMiddleWare())
	router.Use(middlewares.CORSMiddleWare())
	router.NoRoute(noRouteHandler())

	//if config.IsRelease() {
	//    router.Use(loadTLS(config.Address()))
	//}

	userRepository := repository.NewUserRepository(db)
	if err := userRepository.Migrate(); err != nil {
		return nil, fmt.Errorf("user migration: %w", err)
	}

	userService := service.NewUserService(userRepository)
	jwtService, err := utilService.NewJWTService(*config.KeyDir.Value)
	if err != nil {
		return nil, fmt.Errorf("jwt service: %w", err)
	}

	apiV1 := router.Group("/v1")
	SetHealthRoutes(apiV1)
	SetAuthRoutes(apiV1, userService, jwtService, logger)
	SetUsersRoutes(apiV1, userService, logger)
	setupDocRouters(apiV1)

	return router, nil
}

func noRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, model.NewGenericNotFoundResponse(c))
	}
}

func setupDocRouters(rg *gin.RouterGroup) {
	rg.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

//func loadTLS(address string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		middleware := secure.New(secure.Options{
//			SSLRedirect: true,
//			SSLHost:     address,
//		})
//
//		if middleware.Process(c.Writer, c.Request) != nil {
//			return
//		}
//
//		c.Next()
//	}
//}
