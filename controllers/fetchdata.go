package controllers

import "github.com/gin-gonic/gin"

func GetNamespaces(c *gin.Context) {
	c.JSON(200, gin.H{
		"namespaces": "List of Namespaces",
	})
}

func GetPods(c *gin.Context) {
	c.JSON(200, gin.H{
		"pods": "List of Pods",
	})
}
