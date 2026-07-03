package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"transportadora/internal/handlers"
	"transportadora/internal/middleware"
)

func Register(r *gin.Engine, db *gorm.DB) {
	r.Static("/static", "./web/static")

	authHandler := handlers.NewAuthHandler(db)
	homeHandler := handlers.NewHomeHandler(db)
	tripHandler := handlers.NewTripHandler(db)
	historyHandler := handlers.NewHistoryHandler(db)

	clientHandler := handlers.NewClientHandler(db)
	driverHandler := handlers.NewDriverHandler(db)
	tractorHandler := handlers.NewTractorHandler(db)
	trailerHandler := handlers.NewTrailerHandler(db)
	reportHandler := handlers.NewReportHandler(db)

	r.GET("/login", authHandler.LoginForm)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", authHandler.Logout)

	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/", homeHandler.Index)

		trips := protected.Group("/trips")
		{
			trips.GET("", tripHandler.List)
			trips.GET("/new", tripHandler.NewForm)
			trips.GET("/new/bocas", tripHandler.LoadCompartments)
			trips.POST("", tripHandler.Create)
			trips.GET("/:id", tripHandler.Show)
			trips.POST("/:id/finish", tripHandler.Finish)
		}

		clients := protected.Group("/clients")
		{
			clients.GET("", clientHandler.List)
			clients.GET("/new", clientHandler.NewForm)
			clients.POST("", clientHandler.Create)
			clients.GET("/:id/edit", clientHandler.EditForm)
			clients.POST("/:id", clientHandler.Update)
			clients.POST("/:id/delete", clientHandler.Delete)
		}

		drivers := protected.Group("/drivers")
		{
			drivers.GET("", driverHandler.List)
			drivers.GET("/new", driverHandler.NewForm)
			drivers.POST("", driverHandler.Create)
			drivers.GET("/:id/edit", driverHandler.EditForm)
			drivers.POST("/:id", driverHandler.Update)
			drivers.POST("/:id/delete", driverHandler.Delete)
		}

		tractors := protected.Group("/tractors")
		{
			tractors.GET("", tractorHandler.List)
			tractors.GET("/new", tractorHandler.NewForm)
			tractors.POST("", tractorHandler.Create)
			tractors.GET("/:id/edit", tractorHandler.EditForm)
			tractors.POST("/:id", tractorHandler.Update)
			tractors.POST("/:id/delete", tractorHandler.Delete)
		}

		trailers := protected.Group("/trailers")
		{
			trailers.GET("", trailerHandler.List)
			trailers.GET("/new", trailerHandler.NewForm)
			trailers.POST("", trailerHandler.Create)
			trailers.GET("/:id/edit", trailerHandler.EditForm)
			trailers.POST("/:id", trailerHandler.Update)
			trailers.POST("/:id/delete", trailerHandler.Delete)

			trailers.POST("/:id/compartments", trailerHandler.AddCompartment)
			trailers.POST("/:id/compartments/:compId/delete", trailerHandler.DeleteCompartment)
		}

		protected.GET("/history", historyHandler.Index)
		reports := protected.Group("/reports")
		{
			reports.GET("/monthly", reportHandler.Monthly)
			reports.GET("/weekly", reportHandler.Weekly)
		}
	}
}
