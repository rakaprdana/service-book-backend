package cmd

import (
	"fmt"
	"net/http"
	"service-book/common/responses"
	"service-book/config"
	"service-book/constants"
	"service-book/controllers"
	"service-book/databases/seeders"
	modelRepairJob "service-book/domains/models/repair-job"
	modelUser "service-book/domains/models/user"
	"service-book/middleware"
	"service-book/repositories"
	"service-book/routes"
	"service-book/services"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()
		config.Init()

		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}

		time.Local = loc
		err = db.AutoMigrate(
			&modelUser.Role{},
			&modelUser.User{},
			&modelRepairJob.RepairJob{},
			&modelRepairJob.RepairJobImage{},
		)

		if err != nil {
			panic(err)
		}

		seeders.NewSeederRegistry(db).Run()
		repository := repositories.NewRepoRegistry(db)
		service := services.NewServiceRegistry(repository)
		controller := controllers.NewControllerRegistry(service)

		router := gin.Default()
		router.Use(middleware.HandlePanic())
		router.NoRoute(func(ctx *gin.Context) {
			ctx.JSON(http.StatusNotFound, responses.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path :%d", http.StatusNotFound),
			})
		})
		router.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, responses.Response{
				Status:  constants.Success,
				Message: "Welcome to service book",
			})
		})
		router.Use(func(ctx *gin.Context) {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Writer.Header().Set("Access-Control-Allow-Method", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-service-name, x-api-key, x-request-at")
			if ctx.Request.Method == "OPTIONS" {
				ctx.AbortWithStatus(204)
				return
			}
			ctx.Next()
		})
		lmt := tollbooth.NewLimiter(
			config.Config.RateLimiterMaxRequest,
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
			},
		)
		router.Use(middleware.RateLimiter(lmt))

		group := router.Group("api/v1")
		route := routes.NewRouteRegistry(group, controller)
		route.Serve()

		port := fmt.Sprintf(":%d", config.Config.Port)
		router.Run(port)
	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
