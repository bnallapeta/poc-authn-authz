package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	ProviderURL  string
}

func NewOIDCConfig(clientID, clientSecret, redirectURL, providerURL string) *OIDCConfig {
	return &OIDCConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		ProviderURL:  providerURL,
	}
}

func (c *OIDCConfig) GetOAuth2Config(ctx context.Context) (*oauth2.Config, *oidc.IDTokenVerifier, error) {
	provider, err := oidc.NewProvider(ctx, c.ProviderURL)
	if err != nil {
		return nil, nil, err
	}

	oidcConfig := &oidc.Config{
		ClientID: c.ClientID,
	}

	config := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier := provider.Verifier(oidcConfig)

	return config, verifier, nil
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		accessToken := session.Get("access_token")
		if user == nil || accessToken == nil {
			// user is not logged in or access token is not available
			log.Println("No user or access token found in session, redirecting to Keycloak...")
			c.Redirect(http.StatusTemporaryRedirect, "/auth/start")
			c.Abort()
			return
		}
		log.Println("User and access token found in session, proceeding with request...")
		c.Next()
	}
}

func GenerateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)

	// SameSite=Lax is required since without this attribute, cookies cannot be shared between sites
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  expiration,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	return state
}

// Enable CORS in gin.Context
func EnableCors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
