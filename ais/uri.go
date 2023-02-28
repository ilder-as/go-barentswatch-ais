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

// OAuthToken returns the full URL to the OAuth 2.0 token endpoint
func (r URLs) OAuthToken() string {
	return r.OAuthBase + r.TokenEndpoint
}

// AIS returns the full URL to the AIS endpoint
func (r URLs) AIS() string {
	return r.APIBase + r.AISEndpoint
}

// SSEAIS returns the full URL to the AIS endpoint that returns data through Server Sent Events (SSE)
func (r URLs) SSEAIS() string {
	return r.APIBase + r.SSEAISEndpoint
}

// Combined returns the full URL to the Combined endpoint
func (r URLs) Combined() string {
	return r.APIBase + r.CombinedEndpoint
}

// SSECombined returns the full URL to the Combined endpoint that returns data through Server Sent Events (SSE)
func (r URLs) SSECombined() string {
	return r.APIBase + r.SSECombinedEndpoint
}

// LatestCombined returns the full URL to the Latest Combined endpoint
func (r URLs) LatestCombined() string {
	return r.APIBase + r.LatestCombinedEndpoint
}

// LatestAIS returns the full URL to the Latest endpoint
func (r URLs) LatestAIS() string {
	return r.APIBase + r.LatestAISEndpoint
}

// OpenAISArea returns the full URL to the Open AIS Area endpoint
func (r URLs) OpenAISArea() string {
	return r.APIBase + r.OpenAISAreaEndpoint
}
