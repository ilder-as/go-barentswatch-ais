package ais

type config struct {
	ClientID     string
	ClientSecret string
	OAuthScopes  []string
	URLs         URLs
}

type ClientOption func(*config)

func SetResourceURLs(urls URLs) ClientOption {
	return func(conf *config) {
		conf.URLs = urls
	}
}

func SetOAuthScopes(scopes []string) ClientOption {
	return func(conf *config) {
		conf.OAuthScopes = scopes
	}
}
