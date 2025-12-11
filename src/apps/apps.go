
package apps


import (
	"fmt"
	"net/http"
	"findai/src/apps/auth"
	"findai/src/apps/views"
	"findai/src/config"
	"findai/src/database"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
)


func Init() *gin.Engine {

	router := gin.Default()

	// Basic /ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//docs
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	router.GET("/docs", gin.WrapH(middleware.SwaggerUI(opts, nil)))
	router.GET("/swagger.yaml", gin.WrapH(http.FileServer(http.Dir("./docs"))))

	db := database.DB()

	authService := &auth.AuthService{Db: db}
	authViews := &views.AuthViews{AuthService: authService}

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authViews.Register)
		authRoutes.POST("/login", authViews.Login)
	}

	return router
}

func Serve() {
	router := Init()
	router.Run(fmt.Sprintf("0.0.0.0:%d", config.Config.Port))
}
