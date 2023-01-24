package ais

import (
	"context"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

// Client is the main workhorse of the http package. It contains configurations for
// authenticating with Barentswatch' API, and the documented API methods.
//
// A client must be constructed with the NewClient factory function.
type Client struct {
	urls       URLs
	httpClient *http.Client
}

func NewClient(clientId string, clientSecret string) *Client {
	urls := DefaultURLs()
	oauthConfig := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     DefaultURLs().OAuthToken(),
		Scopes:       []string{"ais"},
	}
	httpClient := oauthConfig.Client(context.Background())

	return &Client{
		urls:       urls,
		httpClient: httpClient,
	}
}
