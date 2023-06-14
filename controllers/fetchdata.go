package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetNamespaces(c *gin.Context) {
	session := sessions.Default(c)
	accessToken := session.Get("access_token")

	if accessToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Implement logic to get namespaces using the access token here

	c.JSON(http.StatusOK, gin.H{
		"namespaces": "List of Namespaces",
	})
}

func GetPods(c *gin.Context) {
	session := sessions.Default(c)
	accessToken := session.Get("access_token")

	if accessToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Implement logic to get pods using the access token here

	c.JSON(http.StatusOK, gin.H{
		"pods": "List of Pods",
	})
}
