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
