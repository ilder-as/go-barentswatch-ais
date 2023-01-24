package ais

var defaultUrls URLs = URLs{
	OAuthBase:              "https://id.barentswatch.no",
	TokenEndpoint:          "/connect/token",
	APIBase:                "https://live.ais.barentswatch.no",
	AISEndpoint:            "/v1/ais",
	SSEAISEndpoint:         "/v1/sse/ais",
	CombinedEndpoint:       "/v1/combined",
	SSECombinedEndpoint:    "/v1/sse/combined",
	LatestAISEndpoint:      "/v1/latest/ais",
	LatestCombinedEndpoint: "/v1/latest/combined",
	OpenAISAreaEndpoint:    "/v1/openaisarea",
}

// DefaultURLs returns the documented URLs to the API endpoints.
func DefaultURLs() URLs {
	return defaultUrls
}

type URLs struct {
	// OAuthBase is the base URL for all OAuth requests.
	//
	// Example: https://id.barentswatch.no
	OAuthBase string

	// TokenEndpoint is the relative location of the OAuth Authorization Server token endpoint.
	//
	// Example: /connect/token
	TokenEndpoint string

	// APIBase is the base URL for all API requests.
	//
	// Example: https://live.ais.barentswatch.no
	APIBase string

	// AISEndpoint is the relative location of the AIS API endpoint.
	//
	// Example: /v1/ais
	AISEndpoint string

	// SSEAISEndpoint is the relative location of the Server Side Event AIS API endpoint.
	//
	// Example: /v1/sse/ais
	SSEAISEndpoint string

	// CombinedEndpoint is the relative location of the combined position and static data AIS API endpoint.
	//
	// Example: /v1/combined
	CombinedEndpoint string

	// SSECombinedEndpoint is the relative location of the Server Side Event combined position and static data AIS
	// API endpoint.
	//
	// Example: /v1/sse/combined
	SSECombinedEndpoint string

	// LatestAISEndpoint is the relative location of latest AIS endpoint.
	//
	// Example: /v1/latest/ais
	LatestAISEndpoint string

	// LatestCombinedEndpoint is the relative location of the latest combined position and static data AIS API endpoint.
	//
	// Example: /v1/latest/combined
	LatestCombinedEndpoint string

	// OpenAISAreaEndpoint is the relative location of the Open AIS Area API endpoint.
	//
	// Example: /v1/openaisarea
	OpenAISAreaEndpoint string
}

func (r URLs) OAuthToken() string {
	return r.OAuthBase + r.TokenEndpoint
}

func (r URLs) AIS() string {
	return r.APIBase + r.AISEndpoint
}

func (r URLs) SSEAIS() string {
	return r.APIBase + r.SSEAISEndpoint
}

func (r URLs) Combined() string {
	return r.APIBase + r.CombinedEndpoint
}

func (r URLs) SSECombined() string {
	return r.APIBase + r.SSECombinedEndpoint
}

func (r URLs) LatestCombined() string {
	return r.APIBase + r.LatestCombinedEndpoint
}

func (r URLs) LatestAIS() string {
	return r.APIBase + r.LatestAISEndpoint
}

func (r URLs) OpenAISArea() string {
	return r.APIBase + r.OpenAISAreaEndpoint
}
