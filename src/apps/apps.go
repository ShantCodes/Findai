package apps

import (
	"fmt"
	"net/http"
	"findai/src/apps/views"
	"findai/src/config"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jmoiron/sqlx"
)

func Init(db *sqlx.DB) *gin.Engine {

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

	views.Init(router, db)

	return router
}

func Serve(db *sqlx.DB) {
	router := Init(db)
	router.Run(fmt.Sprintf("0.0.0.0:%d", config.Config.Port))
}
