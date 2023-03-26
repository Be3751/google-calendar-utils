package calendar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

const (
	localAdd           = ":8989"
	redirectUriPattern = "/redirect"
	queryParamName     = "code"
)

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	authCode, err := getAuthCodeQueryParam(localAdd, redirectUriPattern, queryParamName)
	if err != nil {
		log.Fatalf("Unable to get auth code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func getAuthCodeQueryParam(localAdd string, uriPattern string, paramName string) (string, error) {
	var authCode string

	codeChan := make(chan string)
	errChan := make(chan error)
	http.HandleFunc(uriPattern, func(w http.ResponseWriter, r *http.Request) {
		queryVals := r.URL.Query()
		codeChan <- queryVals.Get(paramName)
	})
	go func() {
		err := http.ListenAndServe(localAdd, nil)
		if err != nil {
			errChan <- fmt.Errorf("Unable to launch local server: %w", err)
		}
	}()

	select {
	case authCode = <-codeChan:
		if authCode == "" {
			return "", errors.New("Unable to get auth code")
		}
		return authCode, nil
	case err := <-errChan:
		return "", err
	}
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
