package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bnallapeta/poc-authn-authz/auth"
	"github.com/bnallapeta/poc-authn-authz/routes"
	"golang.org/x/oauth2"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize a session store
	store := cookie.NewStore([]byte("0ya0E6mXjJF9i8IQ"))

	r := gin.Default()

	// Setup gin router to handle CORS
	// This is required as by default, browser disallows js on our page to make a cross-origin request
	// thereby blocking the redirect call to keycloak
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://webapp-route-poc.apps.vmw-sno5.lab.kubeapp.cloud"}
	r.Use(cors.New(corsConfig))

	// Use the session middleware
	r.Use(sessions.Sessions("usersession", store))

	oidcConfig := auth.NewOIDCConfig(
		"poc-webapp",
		"zNszZddURNUokEYIrT8d7mFr9cxacI1p",
		"http://webapp-route-poc.apps.vmw-sno5.lab.kubeapp.cloud/auth/callback",
		"https://keycloak-default.apps.vmw-sno5.lab.kubeapp.cloud/realms/poc-realm",
	)

	oauth2Config, verifier, err := oidcConfig.GetOAuth2Config(context.Background())
	if err != nil {
		log.Fatalf("Failed to get OAuth2 configuration: %v", err)
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/auth/start", func(c *gin.Context) {
		// Enable CORS
		auth.EnableCors(c)

		// Get the path the user was trying to access
		next := c.Request.URL.Query().Get("next")
		if next == "" {
			next = "/"
		}

		// Generate a new OAuth state string for this request
		oAuth2State := auth.GenerateStateOauthCookie(c.Writer)

		// Start a new OAuth2 flow
		authUrl := oauth2Config.AuthCodeURL(oAuth2State, oauth2.AccessTypeOnline)
		log.Printf("Starting new OAuth2 flow with redirect URL: %s", authUrl)

		// Redirect the user to the Keycloak login page
		c.Redirect(http.StatusTemporaryRedirect, authUrl)
	})

	r.GET("/auth/callback", func(c *gin.Context) {
		// Handle callback from Keycloak
		// Parse the authorization code from the request parameters
		// http://webapp-route-poc.apps.vmw-sno5.lab.kubeapp.cloud/auth/callback?next=%2Fapi%2Fgetpods&state=4ZC4PBjUwMYI3IKs&session_state=cc9b44c8-aa66-42dc-add5-0f0054f8798d&code=a3e18706-cac5-4e35-a0b9-afe23fb9bd8c.cc9b44c8-aa66-42dc-add5-0f0054f8798d.b7abba06-1192-49d7-b2b8-ab29d87c2598

		code := c.Request.URL.Query().Get("code")
		if code == "" {
			log.Println("no code found in request")
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"message": "No code found in request",
			})
			return
		}

		// Exchange the authorization code for a token
		log.Printf("Exchanging code: %s with redirect_uri: %s", code, oauth2Config.RedirectURL)
		oAuth2Token, err := oauth2Config.Exchange(context.Background(), code)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"message": fmt.Sprintf("Failed to exchange token: %v", err),
			})
			return
		}

		// Save the access token to session
		session := sessions.Default(c)
		session.Set("access_token", oAuth2Token.AccessToken)
		session.Save()

		// Extract the ID token from OAuth2 token.
		rawIDToken, ok := oAuth2Token.Extra("id_token").(string)
		if !ok {
			log.Println("No id_token field in oauth2 token")
			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"message": "No id_token field in oauth2 token",
			})
			return
		}

		// Parse and verify ID Token payload.
		idToken, err := verifier.Verify(context.Background(), rawIDToken)
		if err != nil {
			log.Printf("Failed to verify ID Token: %v", err)

			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"message": fmt.Sprintf("Failed to verify ID Token: %v", err),
			})
			return
		}

		// Extract custom claims
		var claims struct {
			Email string `json:"email"`
		}
		if err := idToken.Claims(&claims); err != nil {
			log.Printf("Failed to parse claims: %v", err)

			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"message": fmt.Sprintf("Failed to parse claims: %v", err),
			})
			return
		}

		session.Set("user", claims.Email)
		err = session.Save()
		if err != nil {
			log.Printf("Failed to save session: %v", err)

			c.HTML(http.StatusBadRequest, "error.html", gin.H{
				"message": fmt.Sprintf("Failed to save session: %v", err),
			})
			return
		}

		// Redirect to the original page the user was trying to access
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	routes.RegisterRoutes(r)

	// Start the server
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
