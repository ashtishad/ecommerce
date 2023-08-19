package app

import (
	"bytes"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"os"
)

type GoogleAuthHandler struct {
	l *log.Logger
}

var oauth2Config = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_AUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_AUTH_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8000/google-callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

// startGoogleLoginHandler starts google auth process, checks if google auth client id and secret is set
// gets the auth url, and redirects to the callback url
func (gh *GoogleAuthHandler) startGoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	if oauth2Config.ClientID == "" || oauth2Config.ClientSecret == "" {
		gh.l.Printf("invalid google auth client id %s or secret %s", oauth2Config.ClientID, oauth2Config.ClientSecret)
		return
	}
	authURL := oauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

// googleCallbackHandler takes the code from google, generates token, validates it
// gets the client, get's data(only name and email) from google, and generates random for rest of the fields
// sends post request to /user endpoint for user creation or update
// sign_up_option will be set as "google" in database and finally returns the response
func (gh *GoogleAuthHandler) googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		writeResponse(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
		return
	}

	token, err := oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		gh.l.Println("could not convert authorization code into a token")
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	client := oauth2Config.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		gh.l.Println("couldn't get user info from google")
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	defer resp.Body.Close()

	var profile map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
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
		//gh.l.Println(err.Error())
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	createUserResp, err := http.Post("http://localhost:8000/user", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		gh.l.Println("couldn't get createUserResp")
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	defer createUserResp.Body.Close()

	createUserRespBody, err := io.ReadAll(createUserResp.Body)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(createUserResp.StatusCode)
	w.Write(createUserRespBody)
}
