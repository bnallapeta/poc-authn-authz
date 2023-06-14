package routes

import (
	"github.com/bnallapeta/poc-authn-authz/auth"
	"github.com/bnallapeta/poc-authn-authz/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.LoadHTMLGlob("web/static/*")
	r.Static("/static", "./web/static")

	api := r.Group("/api")
	api.Use(auth.AuthRequired())
	{
		api.GET("/getnamespaces", controllers.GetNamespaces)
		api.GET("/getpods", controllers.GetPods)
	}
}
