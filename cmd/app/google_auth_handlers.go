package app

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"os"
)

// GoogleAuthHandler handles Google authentication.
type GoogleAuthHandler struct {
	l *log.Logger
}

var oauth2Config = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_AUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_AUTH_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_AUTH_REDIRECT_URL"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

// StartGoogleLoginHandler starts a Google auth process, checks if google auth client id and secret is set
// gets the auth url, and redirects to the callback url
func (gh *GoogleAuthHandler) StartGoogleLoginHandler(c *gin.Context) {
	if oauth2Config.ClientID == "" || oauth2Config.ClientSecret == "" || oauth2Config.RedirectURL == "" {
		gh.l.Printf("missing %s or secret %s or redirect url %s", oauth2Config.ClientID, oauth2Config.ClientSecret, oauth2Config.RedirectURL)
		return
	}
	authURL := oauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusSeeOther, authURL)
}

// GoogleCallbackHandler takes the code from Google, generates token, validates it
// gets the client, gets data(only name and email) from Google, and generates random for rest of the fields
// sends post request to /user endpoint for user creation or update
// sign_up_option will be set as "google" in database and finally returns the response
func (gh *GoogleAuthHandler) GoogleCallbackHandler(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	token, err := oauth2Config.Exchange(c, code)
	if err != nil {
		gh.l.Println("could not convert authorization code into a token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	client := oauth2Config.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		gh.l.Println("couldn't get user info from google")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer resp.Body.Close()

	var profile map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	jsonData := map[string]string{
		"email":          profile["email"].(string),
		"password":       "securepassword", // example
		"full_name":      profile["name"].(string),
		"phone":          "1234567890", // example
		"sign_up_option": "google",
	}

	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	createUserResp, err := http.Post("http://localhost:8000/users", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		gh.l.Println("couldn't get createUserResp")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer createUserResp.Body.Close()

	createUserRespBody, err := io.ReadAll(createUserResp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Data(createUserResp.StatusCode, "application/json", createUserRespBody)
}
