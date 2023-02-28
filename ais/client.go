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

// NewClient creates a new Client.
//
// It must be called with the user's OAuth client ID and client secret, which can be obtained from Barentswatch.
// Optionally you can supply a single set of URLs to override the default URLs for the API endpoints.
// Do not supply more than zero or one set of URLs.
func NewClient(clientId string, clientSecret string, urls ...URLs) *Client {
	u := DefaultURLs()

	if len(urls) > 0 {
		u = urls[0]
	}
	oauthConfig := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     u.OAuthToken(),
		Scopes:       []string{"ais"},
	}
	httpClient := oauthConfig.Client(context.Background())

	return &Client{
		urls:       u,
		httpClient: httpClient,
	}
}
