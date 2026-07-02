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
	truckHandler := handlers.NewTruckHandler(db)
	tripHandler := handlers.NewTripHandler(db)
	historyHandler := handlers.NewHistoryHandler(db)

	r.GET("/login", authHandler.LoginForm)
	r.POST("/login", authHandler.Login)
	r.POST("/logout", authHandler.Logout)

	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/", homeHandler.Index)

		trucks := protected.Group("/trucks")
		{
			trucks.GET("", truckHandler.List)
			trucks.GET("/new", truckHandler.NewForm)
			trucks.POST("", truckHandler.Create)

			trucks.GET("/:id/edit", truckHandler.EditForm)
			trucks.POST("/:id", truckHandler.Update)
			trucks.POST("/:id/delete", truckHandler.Delete)

			trucks.POST("/:id/compartments", truckHandler.AddCompartment)
			trucks.POST("/:id/compartments/:compId/delete", truckHandler.DeleteCompartment)
		}

		trips := protected.Group("/trips")
		{
			trips.GET("/new", tripHandler.NewForm)
			trips.GET("/new/bocas", tripHandler.LoadCompartments)
			trips.POST("", tripHandler.Create)
			trips.GET("/:id", tripHandler.Show)
			trips.POST("/:id/finish", tripHandler.Finish)
		}

		protected.GET("/history", historyHandler.Index)
	}
}
