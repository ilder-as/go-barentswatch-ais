package ais

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ilder-as/go-barentswatch-ais/ais/option"
	geojson "github.com/paulmach/go.geojson"
	"net/http"
)

// GetAis carries out GET against /v1/ais
func (c *Client) GetAis() (StreamResponse[AisMultiple], error) {
	return c.GetAisContext(context.Background())
}

// GetAisContext carries out GET against /v1/ais with a context for cancellation.
func (c *Client) GetAisContext(ctx context.Context) (StreamResponse[AisMultiple], error) {
	req, err := http.NewRequest("GET", c.urls.AIS(), nil)
	if err != nil {
		return StreamResponse[AisMultiple]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	return StreamResponse[AisMultiple]{Response: res, ctx: ctx, streamType: Simple}, err
}

// PostAis carries out POST against /v1/ais
func (c *Client) PostAis(filterInput FilterInput) (StreamResponse[AisMultiple], error) {
	return c.PostAisContext(context.Background(), filterInput)
}

// PostAisContext carries out POST against /v1/ais with a context for cancellation.
func (c *Client) PostAisContext(ctx context.Context, filterInput FilterInput) (StreamResponse[AisMultiple], error) {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(filterInput); err != nil {
		return StreamResponse[AisMultiple]{}, err
	}
	req, err := http.NewRequest("POST", c.urls.AIS(), body)
	if err != nil {
		return StreamResponse[AisMultiple]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	return StreamResponse[AisMultiple]{Response: res, ctx: ctx, streamType: Simple}, err
}

// GetSSEAis carries out GET against /v1/sse/ais
func (c *Client) GetSSEAis() (StreamResponse[AisMultiple], error) {
	return c.GetSSEAisContext(context.Background())
}

// GetSSEAisContext carries out GET against /v1/sse/ais with a context for cancellation.
func (c *Client) GetSSEAisContext(ctx context.Context) (StreamResponse[AisMultiple], error) {
	req, err := http.NewRequest("GET", c.urls.SSEAIS(), nil)
	if err != nil {
		return StreamResponse[AisMultiple]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	return StreamResponse[AisMultiple]{Response: res, ctx: ctx, streamType: SSE}, err
}

func (c *Client) PostSSEAis(filterInput FilterInput) (StreamResponse[AisMultiple], error) {
	return c.PostSSEAisContext(context.Background(), filterInput)
}

// PostSSEAisContext carries out POST against /v1/sse/ais with a context for cancellation.
func (c *Client) PostSSEAisContext(ctx context.Context, filterInput FilterInput) (StreamResponse[AisMultiple], error) {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(filterInput); err != nil {
		return StreamResponse[AisMultiple]{}, err
	}
	req, err := http.NewRequest("POST", c.urls.AIS(), body)
	if err != nil {
		return StreamResponse[AisMultiple]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	return StreamResponse[AisMultiple]{Response: res, ctx: ctx, streamType: Simple}, err
}

// GetCombined carries out GET against /v1/combined
func (c *Client) GetCombined() (StreamResponse[Combined], error) {
	return c.GetCombinedContext(context.Background())
}

// GetCombinedContext carries out GET against /v1/combined with a context for cancellation.
func (c *Client) GetCombinedContext(ctx context.Context) (StreamResponse[Combined], error) {
	req, err := http.NewRequest("GET", c.urls.Combined(), nil)
	if err != nil {
		return StreamResponse[Combined]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	return StreamResponse[Combined]{Response: res, ctx: ctx, streamType: Simple}, err
}

// PostCombined carries out POST against /v1/combined
func (c *Client) PostCombined(filterInput CombinedFilterInput) (StreamResponse[CombinedMultiple], error) {
	return c.PostCombinedContext(context.Background(), filterInput)
}

// PostCombinedContext carries out POST against /v1/combined with a context for cancellation.
func (c *Client) PostCombinedContext(ctx context.Context, filterInput CombinedFilterInput) (StreamResponse[CombinedMultiple], error) {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(filterInput); err != nil {
		return StreamResponse[CombinedMultiple]{}, err
	}
	req, err := http.NewRequest("POST", c.urls.Combined(), body)
	if err != nil {
		return StreamResponse[CombinedMultiple]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	return StreamResponse[CombinedMultiple]{Response: res, ctx: ctx, streamType: Simple}, err
}

// GetSSECombined carries out GET against /v1/combined
func (c *Client) GetSSECombined() (StreamResponse[Combined], error) {
	return c.GetSSECombinedContext(context.Background())
}

// GetSSECombinedContext carries out GET against /v1/combined with a context for cancellation.
func (c *Client) GetSSECombinedContext(ctx context.Context) (StreamResponse[Combined], error) {
	req, err := http.NewRequest("GET", c.urls.SSECombined(), nil)
	if err != nil {
		return StreamResponse[Combined]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	return StreamResponse[Combined]{Response: res, ctx: ctx, streamType: SSE}, err
}

// PostSSECombined carries out POST against /v1/combined
func (c *Client) PostSSECombined(filterInput CombinedFilterInput) (StreamResponse[CombinedMultiple], error) {
	return c.PostSSECombinedContext(context.Background(), filterInput)
}

// PostSSECombinedContext carries out POST against /v1/combined with a context for cancellation.
func (c *Client) PostSSECombinedContext(ctx context.Context, filterInput CombinedFilterInput) (StreamResponse[CombinedMultiple], error) {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(filterInput); err != nil {
		return StreamResponse[CombinedMultiple]{}, err
	}
	req, err := http.NewRequest("POST", c.urls.SSECombined(), body)
	if err != nil {
		return StreamResponse[CombinedMultiple]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	return StreamResponse[CombinedMultiple]{Response: res, ctx: ctx, streamType: SSE}, err
}

// GetLatestAis carries out GET against /v1/latest/ais
func (c *Client) GetLatestAis(opts ...option.Option) (Response[[]AisMultiple], error) {
	return c.GetLatestAisContext(context.Background(), opts...)
}

// GetLatestAisContext carries out GET against /v1/latest/ais with a context for cancellation.
func (c *Client) GetLatestAisContext(ctx context.Context, opts ...option.Option) (Response[[]AisMultiple], error) {
	req, err := http.NewRequest("GET", c.urls.LatestAIS(), nil)
	if err != nil {
		return Response[[]AisMultiple]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	for _, opt := range opts {
		opt(req)
	}
	res, err := c.httpClient.Do(req)
	return Response[[]AisMultiple]{res}, err
}

// PostLatestAis carries out POST against /v1/latest/ais.
func (c *Client) PostLatestAis(filter LatestAisFilterInput) (Response[[]AisMultiple], error) {
	return c.PostLatestAisContext(context.Background(), filter)
}

// PostLatestAisContext carries out POST against /v1/latest/ais with a context for cancellation.
func (c *Client) PostLatestAisContext(ctx context.Context, filter LatestAisFilterInput) (Response[[]AisMultiple], error) {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(filter); err != nil {
		return Response[[]AisMultiple]{}, err
	}
	req, err := http.NewRequest("POST", c.urls.LatestAIS(), body)
	if err != nil {
		return Response[[]AisMultiple]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	return Response[[]AisMultiple]{res}, err
}

// GetLatestCombined carries out GET against /v1/latest/combined
func (c *Client) GetLatestCombined(opts ...option.Option) (Response[[]Combined], error) {
	return c.GetLatestCombinedContext(context.Background(), opts...)
}

// GetLatestCombinedContext carries out GET against /v1/latest/combined with a context for cancellation.
func (c *Client) GetLatestCombinedContext(ctx context.Context, opts ...option.Option) (Response[[]Combined], error) {
	req, err := http.NewRequest("GET", c.urls.LatestCombined(), nil)
	if err != nil {
		return Response[[]Combined]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	for _, opt := range opts {
		opt(req)
	}
	res, err := c.httpClient.Do(req)
	return Response[[]Combined]{res}, err
}

// GetOpenAisArea carries out GET against /v1/openaisarea with a context for cancellation.
func (c *Client) GetOpenAisArea(ctx context.Context) (Response[geojson.Geometry], error) {
	req, err := http.NewRequest("GET", c.urls.OpenAISArea(), nil)
	if err != nil {
		return Response[geojson.Geometry]{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	return Response[geojson.Geometry]{res}, err
}
