package ais

type config struct {
	ClientID     string
	ClientSecret string
	OAuthScopes  []string
	URLs         URLs
}

type ClientOption func(*config)

func SetResourceURLs(uris URLs) ClientOption {
	return func(conf *config) {
		conf.URLs = uris
	}
}

func SetOAuthScopes(scopes []string) ClientOption {
	return func(conf *config) {
		conf.OAuthScopes = scopes
	}
}
