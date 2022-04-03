package main

import (
	"log"
	"net/http"
	"time"

	"wh2o-next/core/alerts"
	cron "wh2o-next/core/cron"
	"wh2o-next/core/exporter"
	gages "wh2o-next/core/gages"
	"wh2o-next/core/user"
	database "wh2o-next/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func Database(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := gin.Default()
	db := database.InitializeDatabase()

	cron.InitializeCronJobs(db)
	// add db to gin context
	router.Use(Database(db))

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8080", "https://wh2o-api.com"},
		AllowMethods: []string{"PUT", "POST", "GET", "DELETE"},
		AllowHeaders: []string{"Origin"},
		// ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	router.Use(static.Serve("/", static.LocalFile("./client/build", true)))
	// must be a better way to handle direct navigation to react router routes
	// wildcard?
	router.Use(static.Serve("/settings", static.LocalFile("./client/build", true)))

	api := router.Group("/api")
	{

		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"result": "yo",
			})
		})

		api.GET("/gages", gages.HandleGetGages)
		api.GET("/gage-sources/:state", gages.HandleGetGageSources)
		api.POST("/gages", gages.HandleCreateGage)
		api.PUT("/gages", gages.HandleUpdateGage)
		api.DELETE("/gages/:id", gages.HandleDeleteGage)

		api.GET("/alerts", alerts.HandleGetAlerts)
		api.POST("/alerts", alerts.HandleCreateAlert)
		api.PUT("/alerts", alerts.HandleUpdateAlert)
		api.DELETE("/alerts/:id", alerts.HandleDeleteAlert)

		api.GET("/user/:id", user.HandleGetSettings)
		api.PUT("/user", user.HandleUpdateUserSettings)

		api.GET("/export", exporter.ExportAllData)
		api.POST("/import", exporter.ImportData)
	}

	router.Run(":3000")

}
